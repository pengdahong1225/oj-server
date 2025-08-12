package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"oj-server/global"
	"oj-server/module/auth"
	"oj-server/module/configManager"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/src/gateway/internal/define"
	"regexp"
	"strconv"
	"time"

	"oj-server/src/gateway/internal/data"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleUserLogin(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.LoginFrom{})
	if !ret {
		return
	}

	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
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
	resp_data := define.LoginResponse{
		Data:        rpc_resp,
		AccessToken: accessToken,
	}

	resp_data.AccessToken = accessToken
	resp.Data = resp_data
	resp.Message = "登录成功"
	ctx.JSON(http.StatusOK, resp)
}
func HandleUserLoginBySms(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.LoginWithSmsForm{})
	if !ret {
		return
	}

	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 验证码校验
	code, err := data.GetSmsCaptcha(form.Mobile)
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
	resp_data := define.LoginResponse{
		Data:        login_resp,
		AccessToken: accessToken,
	}
	resp.Data = resp_data
	resp.Message = "登录成功"
	ctx.JSON(http.StatusOK, resp)
}

func HandleReFreshAccessToken(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	// 从cookie中获取refresh_token
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		resp.ErrCode = pb.Error_EN_Unauthorized
		resp.Message = "refresh_token不存在"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	j := auth.JWTCreator{
		SigningKey: []byte(configManager.AppConf.JwtCfg.SigningKey),
	}
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, auth.TokenExpired) {
			resp.ErrCode = pb.Error_EN_RefreshTokenExpired
			resp.Message = "refresh_token已过期"
			ctx.JSON(http.StatusUnauthorized, resp)
			return
		} else {
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
	signingKey := configManager.AppConf.JwtCfg.SigningKey
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
			ExpiresAt: time.Now().Unix() + global.RefreshTokenTimeOut, // 7天过期
			Issuer:    global.Issuer,                                  // 签名机构
		},
	}
	return j.CreateToken(claims)
}
func createAccessToken(uid int64, mobile string, role int32) (string, error) {
	signingKey := configManager.AppConf.JwtCfg.SigningKey
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
			ExpiresAt: time.Now().Unix() + global.AccessTokenTimeOut, // 15分钟过期
			Issuer:    global.Issuer,                                 // 签名机构
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
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
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
	form, ret := validate(ctx, define.ResetPasswordForm{})
	if !ret {
		return
	}

	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 验证码校验
	code, err := data.GetSmsCaptcha(form.Mobile)
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
	resp := &define.Response{
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
	rpc_resp, err := client.GetUserInfo(ctx, request)
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

// 处理获取用户某个提交记录的具体信息
func HandleGetUserRecord(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "提交记录id不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 调用题目服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewProblemServiceClient(conn)
	request := &pb.GetSubmitRecordRequest{
		Id: int64(id),
	}
	rpc_resp, err := client.GetSubmitRecordData(ctx, request)
	if err != nil {
		logrus.Errorf("调用题目服务失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "查询失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ErrCode = pb.Error_EN_Success
	resp.Data = rpc_resp.Data
	resp.Message = "获取提交记录成功"
	ctx.JSON(http.StatusOK, resp)
}

// 处理获取用户历史提交记录
// 偏移量分页
func HandleGetUserRecordList(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	// 获取元数据
	uid := ctx.GetInt64("uid")

	// 参数校验
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "页码参数错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "页大小参数错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 调用题目服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer conn.Close()
	client := pb.NewProblemServiceClient(conn)
	req := &pb.GetSubmitRecordListRequest{
		Uid:      uid,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	rpc_resp, err := client.GetSubmitRecordList(context.Background(), req)
	if err != nil {
		logrus.Errorf("problem服务获取提交记录列表失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取提交记录列表失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rpc_resp
	resp.Message = "获取提交记录列表成功"
	resp.ErrCode = pb.Error_EN_Success
	ctx.JSON(http.StatusOK, resp)
}

// 处理获取用户的AC题目列表
func HandleGetUserSolvedList(ctx *gin.Context) {}
