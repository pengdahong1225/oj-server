package handler

import (
	"errors"
	"fmt"
	"net/http"
	"oj-server/global"
	"oj-server/module/configs"
	"oj-server/module/registry"
	"regexp"
	"time"

	"oj-server/svr/gateway/internal/repo"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/middlewares"
	"oj-server/svr/gateway/internal/model"
)

func HandleUserLogin(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.LoginFrom{})
	if !ok {
		return
	}

	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserLoginRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	rpc_resp, err := client.UserLogin(ctx, req)
	if err != nil {
		logrus.Infof("UserLogin Failed: %s", err.Error())
		resp.ErrCode = pb.Error_EN_LoginFailed
		resp.Message = "登录失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	// 生成token
	refreshToken, err := createRefreshAccessToken(rpc_resp.Uid, rpc_resp.Mobile, rpc_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	accessToken, err := createAccessToken(rpc_resp.Uid, rpc_resp.Mobile, rpc_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)
	resp_data := model.LoginResponse{
		UserInfo:    rpc_resp,
		AccessToken: accessToken,
	}

	resp_data.AccessToken = accessToken
	resp.Data = resp_data
	resp.Message = "登录成功"
	ctx.JSON(http.StatusOK, resp)
}
func HandleUserLoginBySms(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.LoginWithSmsForm{})
	if !ok {
		return
	}

	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 验证码校验
	code, err := repo.GetSmsCaptcha(form.Mobile)
	if err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "验证码已过期"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	if form.CaptchaVal != code {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "验证码错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserLoginBySmsCodeRequest{
		Mobile: form.Mobile,
	}
	login_resp, err := client.UserLoginBySmsCode(ctx, req)
	if err != nil {
		logrus.Info("UserLogin Failed: %s", err.Error())
		resp.ErrCode = pb.Error_EN_LoginFailed
		resp.Message = "登录失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	// 生成token
	refreshToken, err := createRefreshAccessToken(login_resp.Uid, login_resp.Mobile, login_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	accessToken, err := createAccessToken(login_resp.Uid, login_resp.Mobile, login_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)
	resp_data := model.LoginResponse{
		UserInfo:    login_resp,
		AccessToken: accessToken,
	}
	resp.Data = resp_data
	resp.Message = "登录成功"
	ctx.JSON(http.StatusOK, resp)
}

func HandleReFreshAccessToken(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}

	// 从cookie中获取refresh_token
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		logrus.Errorf("获取refresh_token失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Unauthorized
		resp.Message = "refresh_token不存在"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	j := middlewares.JWTCreator{
		SigningKey: []byte(configs.AppConf.JwtCfg.SigningKey),
	}
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, middlewares.TokenExpired) {
			logrus.Errorf("refresh_token已过期:%s", err.Error())
			resp.ErrCode = pb.Error_EN_RefreshTokenExpired
			resp.Message = "refresh_token已过期"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		} else {
			logrus.Errorf("refresh_token验证失败:%s", err.Error())
			resp.ErrCode = pb.Error_EN_TokenInvalid
			resp.Message = "refresh_token验证失败"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		}
	}
	// 获取新的access-token
	accessToken, err := createAccessToken(claims.Uid, claims.Mobile, claims.Authority)
	if err != nil {
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "生成token失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Data = gin.H{
		"access_token": accessToken,
	}
	ctx.JSON(http.StatusOK, resp)
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

	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 密码校验
	if form.PassWord != form.RePassWord {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "两次密码输入不匹配"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserRegisterRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	rsp, err := client.UserRegister(ctx, req)
	if err != nil {
		logrus.Info("UserRegister Failed: %s", err.Error())
		resp.ErrCode = pb.Error_EN_RegisterFailed
		resp.Message = "注册失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rsp
	resp.Message = "注册成功"
	ctx.JSON(http.StatusOK, resp)
	return
}
func HandleUserResetPassword(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.ResetPasswordForm{})
	if !ok {
		return
	}

	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 验证码校验
	code, err := repo.GetSmsCaptcha(form.Mobile)
	if err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "验证码已过期"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	if form.CaptchaVal != code {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "验证码错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	req := &pb.ResetUserPasswordRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	rpc_resp, err := client.ResetUserPassword(ctx, req)
	if err != nil {
		logrus.Info("UserLogin Failed: %s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "重置密码失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rpc_resp
	resp.Message = "重置成功"
	ctx.JSON(http.StatusOK, resp)
}

func HandleGetUserProfile(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
	}
	uid := ctx.GetInt64("uid")
	// 调用用户服务
	conn, err := registry.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	request := &pb.GetUserInfoRequest{
		Uid: uid,
	}
	rpc_resp, err := client.GetUserInfoByUid(ctx, request)
	if err != nil {
		logrus.Errorf("用户服务获取用户信息失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "用户信息查询失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rpc_resp
	resp.ErrCode = pb.Error_EN_Success
	resp.Message = "获取用户信息成功"
	ctx.JSON(http.StatusOK, resp)
}
