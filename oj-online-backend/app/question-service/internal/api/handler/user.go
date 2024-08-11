package handler

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/redis"
	"github.com/pengdahong1225/Oj-Online-Server/common/registry"
	"github.com/pengdahong1225/Oj-Online-Server/common/settings"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// 签名过期时间
var (
	jwtTimeOut int64 = 60 * 60 * 24 * 7
	issuer           = "Messi"
)

type UserHandler struct {
}

func (receiver UserHandler) HandleLogin(form *models.LoginFrom) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	mobile, _ := strconv.ParseInt(form.Mobile, 10, 64)
	request := &pb.GetUserDataByMobileRequest{
		Mobile: mobile,
	}
	response, err := client.GetUserDataByMobile(context.Background(), request)
	if err != nil {
		res.Message = err.Error()
		logrus.Debugln(err.Error())
		return res
	}

	// 生成token
	j := middlewares.NewJWT()
	// 设置 payload有效载荷
	claims := &middlewares.UserClaims{
		Uid:       response.Data.Uid,
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

	res.Data = jsonMap
	res.Message = "OK"
	return res
}

func (receiver UserHandler) HandleGetUserProfile(uid int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetUserDataByUid(context.Background(), &pb.GetUserDataByUidRequest{
		Id: uid,
	})
	if err != nil {
		res.Message = err.Error()
		return res
	}

	res.Message = "OK"
	res.Data = response.Data
	return res
}

func (receiver UserHandler) HandleGetRankList() *models.Response {
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

func (receiver UserHandler) HandleGetSubmitRecord(uid int64, stamp int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetUserSubmitRecord(context.Background(), &pb.GetUserSubmitRecordRequest{UserId: uid, Stamp: stamp})
	if err != nil {
		res.Message = err.Error()
		return res
	}
	res.Message = "OK"
	res.Data = response.Data
	return res
}

// HandleGetUserSolvedList 获取用户解决了哪些题目
func (receiver UserHandler) HandleGetUserSolvedList(uid int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetUserSolvedList(context.Background(), &pb.GetUserSolvedListRequest{Uid: uid})
	if err != nil {
		res.Message = err.Error()
		return res
	}
	res.Message = "OK"
	res.Data = response.ProblemSolvedList

	return res
}
