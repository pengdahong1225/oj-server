package api

import (
	"context"
	"errors"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/jwt_utils"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/config"
	"oj-server/svr/gateway/internal/model"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	clients *client.Clients
}

func NewUserHandler(clients *client.Clients) *UserHandler {
	return &UserHandler{
		clients: clients,
	}
}

func (u *UserHandler) HandleUserLogin(ctx *gin.Context) {
	form, ok := validateWithForm(ctx, model.LoginFrom{})
	if !ok {
		return
	}
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.UserLogin(c, &pb.UserLoginRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	})
	if err != nil {
		logrus.Errorf("UserLogin Failed: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	// 生成token
	refreshToken, err := createRefreshAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成refreshToken失败:%s", err.Error())
		ResponseError(ctx, pb.Error_EN_Failed, "生成token失败")
		return
	}
	accessToken, err := createAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成accessToken失败:%s", err.Error())
		ResponseError(ctx, pb.Error_EN_Failed, "生成token失败")
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)

	ResponseOK(ctx, &model.LoginResponse{
		UserInfo:    resp,
		AccessToken: accessToken,
	})
}
func (u *UserHandler) HandleUserLoginBySms(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithForm(ctx, model.LoginWithSmsForm{})
	if !ok {
		return
	}

	// 手机号校验
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.UserLoginBySmsCode(c, &pb.UserLoginBySmsCodeRequest{
		Mobile: form.Mobile,
		Code:   form.CaptchaVal,
	})
	if err != nil {
		logrus.Info("login failed: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	// 生成token
	refreshToken, err := createRefreshAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseError(ctx, pb.Error_EN_Failed, "生成token失败")
		return
	}
	accessToken, err := createAccessToken(resp.Uid, resp.Mobile, resp.Role)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseError(ctx, pb.Error_EN_Failed, "生成token失败")
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)

	ResponseOK(ctx, &model.LoginResponse{
		UserInfo:    resp,
		AccessToken: accessToken,
	})
}

func (u *UserHandler) HandleReFreshAccessToken(ctx *gin.Context) {
	// 从cookie中获取refresh_token
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		logrus.Errorf("获取refresh_token失败:%s", err.Error())
		ResponseBadRequest(ctx, "refresh_token不存在")
		return
	}
	// 解析refresh_token
	jwt_builder := jwt_utils.JWTBuilder{
		SigningKey: []byte(config.MyConfig.JwtCfg.SigningKey),
	}
	claims := new(model.UserClaims)
	if err = jwt_builder.ParseToken(refreshToken, claims); err != nil {
		if errors.Is(err, jwt_utils.TokenExpired) {
			logrus.Errorf("refresh_token已过期:%s", err.Error())
			ResponseUnauthorized(ctx, "refresh_token已过期")
			return
		} else {
			logrus.Errorf("refresh_token验证失败:%s", err.Error())
			ResponseUnauthorized(ctx, "refresh_token验证失败")
			return
		}
	}
	// 获取新的access-token
	accessToken, err := createAccessToken(claims.Uid, claims.Mobile, claims.Authority)
	if err != nil {
		logrus.Errorf("生成token失败:%s", err.Error())
		ResponseError(ctx, pb.Error_EN_Failed, "生成token失败")
		return
	}

	ResponseOK(ctx, gin.H{
		"access_token": accessToken,
	})
}

func (u *UserHandler) HandleUserRegister(ctx *gin.Context) {
	form, ok := validateWithForm(ctx, model.RegisterForm{})
	if !ok {
		return
	}
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}
	if form.PassWord != form.RePassWord {
		ResponseBadRequest(ctx, "两次密码输入不匹配")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.UserRegister(c, &pb.UserRegisterRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
	})
	if err != nil {
		logrus.Errorf("register failed: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gin.H{
		"uid": resp.Uid,
	})
}
func (u *UserHandler) HandleUserResetPassword(ctx *gin.Context) {
	form, ok := validateWithForm(ctx, model.ResetPasswordForm{})
	if !ok {
		return
	}
	ok, _ = regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := u.clients.User.ResetUserPassword(c, &pb.ResetUserPasswordRequest{
		Mobile:   form.Mobile,
		Password: form.PassWord,
		Code:     form.CaptchaVal,
	})
	if err != nil {
		logrus.Info("ResetUserPassword Failed: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

func (u *UserHandler) HandleGetUserProfile(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.GetUserInfoByUid(c, &pb.GetUserInfoRequest{Uid: uid})
	if err != nil {
		logrus.Errorf("用户服务获取用户信息失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, resp.Data)
}

func (u *UserHandler) HandleGetUserList(ctx *gin.Context) {
	var params model.QueryUserListParams
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		logrus.Errorf("参数校验失败: %s", err.Error())
		ResponseBadRequest(ctx, "参数验证失败")
		return
	}
	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.GetUserList(c, &pb.GetUserListRequest{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
	})
	if err != nil {
		logrus.Errorf("用户服务获取用户列表失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	result := &model.QueryUserListResult{
		Total: resp.Total,
	}
	for _, v := range resp.Data {
		result.List = append(result.List, &model.UserInfo{
			Uid:      v.Uid,
			Mobile:   v.Mobile,
			Nickname: v.Nickname,
			Email:    v.Email,
			Gender:   v.Gender,
			CreateAt: v.CreateAt,
			Role:     v.Role,
		})
	}

	ResponseOK(ctx, result)
}
func createRefreshAccessToken(uid int64, mobile string, role int32) (string, error) {
	jwt_builder := jwt_utils.JWTBuilder{
		SigningKey: []byte(config.MyConfig.JwtCfg.SigningKey),
	}
	claims := &model.UserClaims{
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
	return jwt_builder.CreateToken(claims)
}
func createAccessToken(uid int64, mobile string, role int32) (string, error) {
	jwt_builder := jwt_utils.JWTBuilder{
		SigningKey: []byte(config.MyConfig.JwtCfg.SigningKey),
	}
	claims := &model.UserClaims{
		Uid:       uid,
		Mobile:    mobile,
		Authority: role,
		Type:      "access",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                             // 签名生效时间
			ExpiresAt: time.Now().Unix() + consts.AccessTokenTimeOut, // 过期时间
			Issuer:    consts.Issuer,                                 // 签名机构
		},
	}
	return jwt_builder.CreateToken(claims)
}
