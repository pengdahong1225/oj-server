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
	req := &model.LoginRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	ok, err := regexp.MatchString(`^1[3-9]\d{9}$`, req.Mobile)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.UserLogin(c, &pb.UserLoginRequest{
		Mobile:   req.Mobile,
		Password: req.PassWord,
	})
	if err != nil {
		logrus.Errorf("user login failed: %s", err.Error())
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
	req := &model.LoginWithSmsRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	ok, err := regexp.MatchString(`^1[3-9]\d{9}$`, req.Mobile)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := u.clients.User.UserLoginBySmsCode(c, &pb.UserLoginBySmsCodeRequest{
		Mobile: req.Mobile,
		Code:   req.CaptchaVal,
	})
	if err != nil {
		logrus.Errorf("user login failed: %s", err.Error())
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
		SigningKey: []byte(config.MyConfig.Jwt.SigningKey),
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

// 退出登录：清除 refresh_token cookie
func (u *UserHandler) HandleUserLogout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)
	ResponseOK(ctx, nil)
}

func (u *UserHandler) HandleUserRegister(ctx *gin.Context) {
	req := &model.RegisterRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	ok, err := regexp.MatchString(`^1[3-9]\d{9}$`, req.Mobile)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := u.clients.User.UserRegister(c, &pb.UserRegisterRequest{
		Mobile:     req.Mobile,
		Password:   req.PassWord,
		RePassword: req.RePassWord,
	})
	if err != nil {
		logrus.Errorf("register failed: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}
func (u *UserHandler) HandleUserResetPassword(ctx *gin.Context) {
	req := &model.ResetPasswordRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	ok, err := regexp.MatchString(`^1[3-9]\d{9}$`, req.Mobile)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	if !ok {
		ResponseBadRequest(ctx, "手机号格式错误")
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = u.clients.User.ResetUserPassword(c, &pb.ResetUserPasswordRequest{
		Mobile:   req.Mobile,
		Password: req.PassWord,
		Code:     req.CaptchaVal,
	})
	if err != nil {
		logrus.Info("reset password failed: %s", err.Error())
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
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, resp.Data)
}

func (u *UserHandler) HandleGetUserList(ctx *gin.Context) {
	req := &model.QueryUserListRequest{}
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	// 调用用户服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := u.clients.User.GetUserList(c, &pb.GetUserListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		logrus.Errorf("用户服务获取用户列表失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}
func createRefreshAccessToken(uid int64, mobile string, role int32) (string, error) {
	jwt_builder := jwt_utils.JWTBuilder{
		SigningKey: []byte(config.MyConfig.Jwt.SigningKey),
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
		SigningKey: []byte(config.MyConfig.Jwt.SigningKey),
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

func (u *UserHandler) HandleGetImageCapatcha(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := u.clients.User.GetImageCapcha(c, nil)
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}
func (u *UserHandler) HandleSendSmsCapatcha(ctx *gin.Context) {
	req := &model.SendSmsCodeRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := u.clients.User.SendSmsCode(c, &pb.SendSmsCodeRequest{
		Mobile:      req.Mobile,
		CaptchaId:   req.CaptchaID,
		CaptchaCode: req.CaptchaCode,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}
