package handler

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/auth"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"time"
)

func HandleUserLogin(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.LoginFrom{})
	if !ret {
		return
	}

	resp := &define.Response{
		Code:    define.Success,
		Message: "",
		Data:    nil,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.Code = define.Failed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(consts.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	req := &pb.UserLoginRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	}
	login_resp, err := client.UserLogin(ctx, req)
	if err != nil {
		logrus.Info("UserLogin Failed: %s", err.Error())
		resp.Code = define.Failed
		resp.Message = "登录失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	data := define.LoginRspData{
		Rsp: login_resp,
	}
	// 生成token
	refreshToken, err := createRefreshAccessToken(login_resp.Uid, login_resp.Mobile, login_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	accessToken, err := createAccessToken(login_resp.Uid, login_resp.Mobile, login_resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)
	data.AccessToken = accessToken
	resp.Data = data
	resp.Message = "登录成功"
	ctx.JSON(http.StatusOK, resp)
	return
}
func HandleReFreshAccessToken(ctx *gin.Context) {
	resp := &define.Response{
		Code: define.Success,
	}

	// 从cookie中获取refresh_token
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		resp.Code = define.Unauthorized
		resp.Message = "refresh_token不存在"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	j := auth.JWTCreator{
		SigningKey: []byte(settings.Instance().SigningKey),
	}
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, auth.TokenExpired) {
			resp.Code = define.RefreshTokenExpired
			resp.Message = "refresh_token已过期"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		} else {
			resp.Code = define.TokenInvalid
			resp.Message = "token验证失败"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		}
	}
	// 获取新的access-token
	accessToken, err := createAccessToken(claims.Uid, claims.Mobile, claims.Authority)
	if err != nil {
		resp.Code = define.Failed
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
	signingKey := settings.Instance().SigningKey
	j := auth.JWTCreator{
		SigningKey: []byte(signingKey),
	}
	claims := &auth.UserClaims{
		Uid:       uid,
		Mobile:    mobile,
		Authority: role,
		Type:      "refresh",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                              // 签名生效时间
			ExpiresAt: time.Now().Unix() + consts.RefreshTokenTimeOut, // 7天过期
			Issuer:    consts.Issuer,                                  // 签名机构
		},
	}
	return j.CreateToken(claims)
}
func createAccessToken(uid int64, mobile string, role int32) (string, error) {
	signingKey := settings.Instance().SigningKey
	j := auth.JWTCreator{
		SigningKey: []byte(signingKey),
	}
	claims := &auth.UserClaims{
		Uid:       uid,
		Mobile:    mobile,
		Authority: role,
		Type:      "access",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                             // 签名生效时间
			ExpiresAt: time.Now().Unix() + consts.AccessTokenTimeOut, // 15分钟过期
			Issuer:    consts.Issuer,                                 // 签名机构
		},
	}
	return j.CreateToken(claims)
}

func HandleUserRegister(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.RegisterForm{})
	if !ret {
		return
	}

	resp := &define.Response{
		Code: define.Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.Code = define.Failed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 密码校验
	if form.PassWord != form.RePassWord {
		resp.Code = define.Failed
		resp.Message = "两次密码输入不匹配"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用用户服务
	conn, err := registry.GetGrpcConnection(consts.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = "服务器错误"
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
		resp.Code = define.Failed
		resp.Message = "注册失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rsp
	resp.Message = "注册成功"
	ctx.JSON(http.StatusOK, resp)
	return
}
func HandleUserResetPassword(ctx *gin.Context) {}
func HandleGetUserProfile(ctx *gin.Context)    {}
func HandleGetUserRecordList(ctx *gin.Context) {}
func HandleGetUserRecord(ctx *gin.Context)     {}
func HandleGetUserSolvedList(ctx *gin.Context) {}
