package captcha

import (
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
	"oj-server/svr/user/internal/configs"
	"oj-server/utils"
	"os"
)

func SendSmsCode(mobile string) (string, error) {
	logMode := os.Getenv("LOG_MODE")
	if logMode == "release" {
		// 生成随机数
		c := utils.GenerateSmsCode(6)
		param := map[string]string{
			"code": c,
		}
		data, _ := json.Marshal(param)

		// 调用第三方服务发送
		if err := send(data, mobile); err != nil {
			return "", err
		} else {
			return c, nil
		}
	}
	return "123456", nil
}

func send(param []byte, phone string) error {
	sms_cfg := configs.AppConf.SmsCfg

	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &sms_cfg.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &sms_cfg.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(sms_cfg.Endpoint)
	client, _ := dysmsapi.NewClient(config)

	request := &dysmsapi.SendSmsRequest{}
	request.SetSignName(sms_cfg.SignName)
	request.SetTemplateCode(sms_cfg.TemplateCode)
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
