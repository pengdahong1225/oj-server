package handler

import (
	"errors"
	"fmt"
	"net/http"
	"oj-server/global"
	"oj-server/pkg/registry"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/configs"
	"oj-server/svr/gateway/internal/middlewares"
	"oj-server/svr/gateway/internal/model"
)

func HandleUserLogin(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.LoginFrom{})
	if !ok {
		return
	}

	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseWithJson(ctx, http.StatusBadRequest, "手机号格式错误", nil)
		return
	}

	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserLoginRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	resp, err := client.UserLogin(ctx, req)
	if err != nil {
		logrus.Infof("UserLogin Failed: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "登录失败", nil)
		return
	}
	// 生成token
	refreshToken, err := createRefreshAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, fmt.Sprintf("生成token失败:%s", err.Error()), nil)
		return
	}
	accessToken, err := createAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, fmt.Sprintf("生成token失败:%s", err.Error()), nil)
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)

	ResponseWithJson(ctx, http.StatusOK, "success", &model.LoginResponse{
		UserInfo:    resp,
		AccessToken: accessToken,
	})
}
func HandleUserLoginBySms(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.LoginWithSmsForm{})
	if !ok {
		return
	}

	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseWithJson(ctx, http.StatusBadRequest, "手机号格式错误", nil)
		return
	}

	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserLoginBySmsCodeRequest{
		Mobile: form.Mobile,
		Code:   form.CaptchaVal,
	}
	resp, err := client.UserLoginBySmsCode(ctx, req)
	if err != nil {
		logrus.Info("UserLogin Failed: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "登录失败", nil)
		return
	}
	// 生成token
	refreshToken, err := createRefreshAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, fmt.Sprintf("生成token失败:%s", err.Error()), nil)
		return
	}
	accessToken, err := createAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, fmt.Sprintf("生成token失败:%s", err.Error()), nil)
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)

	ResponseWithJson(ctx, http.StatusOK, "success", &model.LoginResponse{
		UserInfo:    resp,
		AccessToken: accessToken,
	})
}

func HandleReFreshAccessToken(ctx *gin.Context) {
	// 从cookie中获取refresh_token
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		logrus.Errorf("获取refresh_token失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusBadRequest, "refresh_token不存在", nil)
		return
	}
	j := middlewares.JWTCreator{
		SigningKey: []byte(configs.AppConf.JwtCfg.SigningKey),
	}
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, middlewares.TokenExpired) {
			logrus.Errorf("refresh_token已过期:%s", err.Error())
			ResponseWithJson(ctx, http.StatusUnauthorized, "refresh_token已过期", nil)
			return
		} else {
			logrus.Errorf("refresh_token验证失败:%s", err.Error())
			ResponseWithJson(ctx, http.StatusUnauthorized, "refresh_token验证失败", nil)
			return
		}
	}
	// 获取新的access-token
	accessToken, err := createAccessToken(claims.Uid, claims.Mobile, claims.Authority)
	if err != nil {
		ResponseWithJson(ctx, http.StatusInternalServerError, "生成token失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", gin.H{
		"access_token": accessToken,
	})
}
func createRefreshAccessToken(uid int64, mobile string, role int32) (string, error) {
	signingKey := configs.AppConf.JwtCfg.SigningKey
	j := middlewares.JWTCreator{
		SigningKey: []byte(signingKey),
	}
	claims := &middlewares.UserClaims{
		Uid:       uid,
		Mobile:    mobile,
		Authority: role,
		Type:      "refresh",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                              // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.RefreshTokenTimeOut, // 7天过期
			Issuer:    global.Issuer,                                  // 签名机构
		},
	}
	return j.CreateToken(claims)
}
func createAccessToken(uid int64, mobile string, role int32) (string, error) {
	signingKey := configs.AppConf.JwtCfg.SigningKey
	j := middlewares.JWTCreator{
		SigningKey: []byte(signingKey),
	}
	claims := &middlewares.UserClaims{
		Uid:       uid,
		Mobile:    mobile,
		Authority: role,
		Type:      "access",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                             // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.AccessTokenTimeOut, // 15分钟过期
			Issuer:    global.Issuer,                                 // 签名机构
		},
	}
	return j.CreateToken(claims)
}

func HandleUserRegister(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.RegisterForm{})
	if !ok {
		return
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseWithJson(ctx, http.StatusBadRequest, "手机号格式错误", nil)
		return
	}
	// 密码校验
	if form.PassWord != form.RePassWord {
		ResponseWithJson(ctx, http.StatusBadRequest, "两次密码输入不匹配", nil)
		return
	}

	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserRegisterRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	rsp, err := client.UserRegister(ctx, req)
	if err != nil {
		logrus.Errorf("UserRegister Failed: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "注册失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", gin.H{
		"uid": rsp.Uid,
	})
}
func HandleUserResetPassword(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.ResetPasswordForm{})
	if !ok {
		return
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseWithJson(ctx, http.StatusBadRequest, "手机号格式错误", nil)
		return
	}

	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewUserServiceClient(conn)
	req := &pb.ResetUserPasswordRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
		Code:     form.CaptchaVal,
	}
	_, err = client.ResetUserPassword(ctx, req)
	if err != nil {
		logrus.Info("ResetUserPassword Failed: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "重置密码失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusInternalServerError, "success", nil)
}
func HandleGetUserProfile(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewUserServiceClient(conn)
	resp, err := client.GetUserInfoByUid(ctx, &pb.GetUserInfoRequest{
		Uid: uid,
	})
	if err != nil {
		logrus.Errorf("用户服务获取用户信息失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "用户信息查询失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", resp.Data)
}
