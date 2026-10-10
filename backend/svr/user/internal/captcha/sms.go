package captcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"oj-server/svr/user/internal/config"
	"os"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

func SendSmsCode(mobile string) (string, error) {
	logMode := os.Getenv("LOG_MODE")
	if logMode == "release" {
		c := generateSmsCode(6)
		param := map[string]string{
			"code": c,
		}
		data, _ := json.Marshal(param)

		if err := send(data, mobile); err != nil {
			return "", err
		} else {
			return c, nil
		}
	}
	return "123456", nil
}

func generateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func send(param []byte, phone string) error {
	sms_cfg := config.MyConfig.Sms

	config := &openapi.Config{
		AccessKeyId:     &sms_cfg.AccessKeyId,
		AccessKeySecret: &sms_cfg.AccessKeySecret,
	}
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
