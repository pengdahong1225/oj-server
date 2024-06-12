package captcha

import (
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
	"os"
	"question-service/services/redis"
	"question-service/settings"
	"question-service/utils"
)

func SendSmsCode(mobile string) error {
	logMode := os.Getenv("LOG_MODE")
	if logMode == "release" {
		// 生成随机数
		c := utils.GenerateSmsCode(6)
		expire := 180 // 3min过期
		param := map[string]string{
			"code": c,
		}
		data, _ := json.Marshal(param)

		// 调用第三方服务发送
		if err := send(data, mobile); err != nil {
			return err
		}

		// 缓存验证码
		if err := redis.SetKVByStringWithExpire(mobile, c, expire); err != nil {
			return err
		}
	} else {
		return nil
	}

	return nil
}

func VerifySmsCode(mobile string, code string) bool {
	logMode := os.Getenv("LOG_MODE")
	if logMode == "release" {
		value, err := redis.GetValueByString(mobile)
		if err != nil {
			logrus.Infoln("redis get value err:", err)
			return false
		} else {
			return value == code
		}
	} else {
		return code == "123456"
	}
}

func send(param []byte, phone string) error {
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
	request.SetTemplateParam(string(param))

	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if *response.Body.Code != "OK" {
		return errors.New(tea.StringValue(response.Body.Message))
	}
	logrus.Debugln(tea.StringValue(response.Body.RequestId))

	return nil
}