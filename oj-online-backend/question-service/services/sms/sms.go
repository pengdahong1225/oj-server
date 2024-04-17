package sms

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
	"question-service/settings"
)

func Send(param []byte, phone string) error {
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
