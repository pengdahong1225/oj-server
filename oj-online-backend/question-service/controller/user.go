package controller

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/logic"
	"question-service/services/redis"
	"question-service/settings"
	"question-service/utils"
	"regexp"
	"strconv"
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
	// 查询参数
	if phone, ok := ctx.GetQuery("phone"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, phone)
		if !ok {
			ctx.Abort()
			return
		}
		phoneInt, _ := strconv.ParseInt(phone, 10, 64)
		logic.GetUserDetail(ctx, phoneInt)
	}
}

func GetRankList(ctx *gin.Context) {
	logic.GetRankList(ctx)
}

func GetSubmitRecord(ctx *gin.Context) {
	if id, ok := ctx.GetQuery("userId"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		userId, _ := strconv.ParseInt(id, 10, 64)
		logic.GetSubmitRecord(ctx, userId)
	}
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
		AccessKeyId: &settings.Conf.SmsConfig.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &settings.Conf.SmsConfig.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(settings.Conf.SmsConfig.Endpoint)
	client, _ := dysmsapi.NewClient(config)

	request := &dysmsapi.SendSmsRequest{}
	request.SetSignName(settings.Conf.SmsConfig.SignName)
	request.SetTemplateCode(settings.Conf.SmsConfig.TemplateCode)
	request.SetPhoneNumbers(phone)
	request.SetTemplateParam(string(data))

	response, errRsp := client.SendSms(request)
	if errRsp != nil || *response.Body.Code != "OK" {
		ctx.JSON(http.StatusInternalServerError, errRsp)
		return
	}
	logrus.Debugln(tea.StringValue(response.Body.RequestId))

	// 缓存验证码
	redisConn := redis.NewConn()
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
