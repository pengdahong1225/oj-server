package controller

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/global"
	"question-service/logic"
	"question-service/utils"
)

func UserRegister(ctx *gin.Context) {
	// 表单验证
	if form, ok := formValidateForRegistry(ctx); ok {
		logic.RegistryHandler(ctx, form)
	}
}

func UserLogin(ctx *gin.Context) {
	// 表单验证
	if form, ok := formValidateForLogin(ctx); ok {
		logic.LoginHandler(ctx, form)
	}
}

func GetUserDetail(ctx *gin.Context) {
	logic.GetUserDetail(ctx)
}

func GetRankList(ctx *gin.Context) {
	logic.GetRankList(ctx)
}

func GetSubmitRecord(ctx *gin.Context) {
	logic.GetSubmitRecord(ctx)
}

func SendCmsCode(ctx *gin.Context) {
	// 生成随机数
	c := utils.GenerateSmsCode(6)
	param := map[string]string{
		"code": c,
	}
	data, _ := json.Marshal(param)

	phone := "18048155008"
	expire := 180 // 3min过期

	// controller
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &global.ConfigInstance.SMS_.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &global.ConfigInstance.SMS_.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(global.ConfigInstance.SMS_.Endpoint)
	client, _ := dysmsapi.NewClient(config)

	request := &dysmsapi.SendSmsRequest{}
	request.SetSignName(global.ConfigInstance.SMS_.SignName)
	request.SetTemplateCode(global.ConfigInstance.SMS_.TemplateCode)
	request.SetPhoneNumbers(phone)
	request.SetTemplateParam(string(data))

	response, errRsp := client.SendSms(request)
	if errRsp != nil || *response.Body.Code != "OK" {
		ctx.JSON(http.StatusInternalServerError, errRsp)
		return
	}
	logrus.Debugln(tea.StringValue(response.Body.RequestId))

	// 缓存验证码
	redisConn := global.RedisPoolInstance.Get()
	defer redisConn.Close()
	if _, err := redisConn.Do("Set", phone, c, "ex", expire); err != nil {
		logrus.Errorln(err)
		ctx.JSON(http.StatusInternalServerError, errRsp)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
