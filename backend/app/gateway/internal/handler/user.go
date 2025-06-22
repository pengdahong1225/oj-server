package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/services"
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
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    define.Failed,
			"message": "手机号格式错误",
		})
		return
	}

	resp := &define.Response{
		Code:    define.Success,
		Message: "",
		Data:    nil,
	}
	defer ctx.JSON(http.StatusOK, resp)

	// 调用用户服务
	conn, err := services.Instance.GetConnection(consts.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = "用户服务连接失败"
		return
	}
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
		return
	}
	data := define.LoginRspData{
		Rsp: login_resp,
	}
	// 生成token
	j := middlewares.NewJWT()
	// 设置 payload有效载荷
	claims := &middlewares.UserClaims{
		Uid:       login_resp.Uid,
		Mobile:    login_resp.Mobile,
		Authority: login_resp.Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                       // 签名生效时间
			ExpiresAt: time.Now().Unix() + consts.TokenTimeOut, // 7天过期
			Issuer:    consts.Issuer,                           // 签名机构
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		resp.Code = define.Failed
		resp.Message = fmt.Sprintf("生成token失败:%s", err.Error())
		resp.Data = nil
		logrus.Errorf("生成token失败:%s", err.Error())
		return
	}
	data.Token = token

	resp.Data = data
	resp.Message = "登录成功"
	return
}
func HandleUserRegister(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.RegisterForm{})
	if !ret {
		return
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    define.Failed,
			"message": "手机号格式错误",
		})
		return
	}
	// 密码校验
	if form.PassWord != form.RePassWord {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    define.Failed,
			"message": "两次密码输入不匹配",
		})
		return
	}

	resp := &define.Response{
		Code:    define.Success,
		Message: "",
		Data:    nil,
	}
	defer ctx.JSON(http.StatusOK, resp)

	// 调用用户服务
	conn, err := services.Instance.GetConnection(consts.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.Code = define.Failed
		resp.Message = "用户服务连接失败"
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
		return
	}
	resp.Data = rsp
	resp.Message = "注册成功"
	return
}
func HandleUserResetPassword(ctx *gin.Context) {

}
func HandleGetUserProfile(ctx *gin.Context) {}

func HandleGetUserRecordList(ctx *gin.Context) {}
func HandleGetUserRecord(ctx *gin.Context)     {}
func HandleGetUserSolvedList(ctx *gin.Context) {}
