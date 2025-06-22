package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
	"net/http"
	"regexp"
	"github.com/sirupsen/logrus"
	"github.com/pengdahong1225/oj-server/backend/module/services"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
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
	rsp, err := client.UserLogin(ctx, req)
	if err != nil {
		logrus.Info("UserLogin Failed: %s", err.Error())
		resp.Code = define.Failed
		resp.Message = "登录失败"
		return
	}
	resp.Data = rsp
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
