package internal

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/api/proto"
	"question-service/middlewares"
	"question-service/models"
	"question-service/services/redis"
	"question-service/services/registry"
	"question-service/services/sms"
	"question-service/settings"
	"question-service/utils"
	"strconv"
	"time"
)

// 签名过期时间
var (
	jwtTimeOut int64 = 60 * 60 * 24 * 7
	issuer           = "Messi"
)

func ProcessForLogin(form *models.LoginFrom) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	mobile, _ := strconv.ParseInt(form.Mobile, 10, 64)
	request := &pb.GetUserRequest{
		Mobile: mobile,
	}
	response, err := client.GetUserData(context.Background(), request)
	if err != nil {
		res.Message = err.Error()
		logrus.Debugln(err.Error())
		return res
	}

	// 生成token
	j := middlewares.NewJWT()
	// 设置 payload有效载荷
	claims := &middlewares.UserClaims{
		Mobile:    response.Data.Mobile,
		Authority: response.Data.Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              // 签名生效时间
			ExpiresAt: time.Now().Unix() + jwtTimeOut, // 7天过期
			Issuer:    issuer,                         // 签名机构
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("生成token失败:%s\n", err.Error())
		return res
	}

	jsonMap := make(map[string]interface{})
	jsonMap["userinfo"] = response.Data
	jsonMap["token"] = token
	bytes, _ := json.Marshal(jsonMap)

	res.Data = string(bytes)
	return res
}

func GetUserDetail(mobile int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetUserData(context.Background(), &pb.GetUserRequest{
		Mobile: mobile,
	})
	if err != nil {
		res.Message = err.Error()
		return res
	}

	res.Message = "OK"
	res.Data = response.Data
	return res
}

func GetRankList() *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	// 获取排行榜
	reply, err := redis.QueryRankList()
	if err != nil {
		res.Message = "排行榜获取失败"
		logrus.Debugf("排行榜获取失败:%s\n", err.Error())
		return res
	}
	var rankList []models.RankList
	for i := 0; i < len(reply); i += 2 {
		user := models.UserInfo{}
		json.Unmarshal([]byte(reply[i]), &user)
		item := models.RankList{
			Phone:     user.Phone,
			NickName:  user.NickName,
			PassCount: user.PassCount,
		}
		rankList = append(rankList, item)
	}

	bytes, err := json.Marshal(rankList)
	if err != nil {
		res.Message = err.Error()
		logrus.Errorf("排行榜序列化失败:%s", err.Error())
		return res
	}

	res.Message = "OK"
	res.Data = string(bytes)
	return res
}

func GetSubmitRecord(userId int64) {

}

func SendSmsCode(mobile string) error {
	// 生成随机数
	c := utils.GenerateSmsCode(6)
	expire := 180 // 3min过期
	param := map[string]string{
		"code": c,
	}
	data, _ := json.Marshal(param)

	// 调用第三方服务发送
	if err := sms.Send(data, mobile); err != nil {
		return err
	}

	// 缓存验证码
	if err := redis.SetKVByStringWithExpire(mobile, c, expire); err != nil {
		return err
	}

	return nil
}