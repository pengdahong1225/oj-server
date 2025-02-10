package logic

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type User struct {
}

func (r User) OnUserRegister(form *models.RegisterForm) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	mobile, _ := strconv.ParseInt(form.Mobile, 10, 64)
	// todo 密码加密
	hash := utils.HashPassword(form.PassWord)

	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewUserServiceClient(dbConn)
	request := &pb.CreateUserRequest{Data: &pb.UserInfo{
		Mobile:   mobile,
		Password: hash,
	}}
	response, err := client.CreateUserData(context.Background(), request)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Debugln(err.Error())
		return res
	}
	res.Message = "ok"
	res.Data = map[string]int64{
		"id": response.Id,
	}
	return res
}

func (r User) OnUserLogin(form *models.LoginFrom) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	mobile, _ := strconv.ParseInt(form.Mobile, 10, 64)

	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewUserServiceClient(dbConn)
	request := &pb.GetUserDataByMobileRequest{
		Mobile: mobile,
	}
	response, err := client.GetUserDataByMobile(context.Background(), request)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Debugln(err.Error())
		return res
	}

	// 校验密码
	if utils.HashPassword(form.PassWord) != response.Data.Password {
		res.Code = models.Failed
		res.Message = "密码错误"
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
			NotBefore: time.Now().Unix(),                       // 签名生效时间
			ExpiresAt: time.Now().Unix() + consts.TokenTimeOut, // 7天过期
			Issuer:    consts.Issuer,                           // 签名机构
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("生成token失败:%s\n", err.Error())
		return res
	}

	res.Data = models.LoginRspData{
		Uid:       response.Data.Uid,
		Mobile:    response.Data.Mobile,
		NickName:  response.Data.Nickname,
		Email:     response.Data.Email,
		Gender:    response.Data.Gender,
		Role:      response.Data.Role,
		AvatarUrl: response.Data.AvatarUrl,
		Token:     token,
	}
	res.Message = "OK"
	return res
}

func (r User) GetUserProfile(uid int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewUserServiceClient(dbConn)
	response, err := client.GetUserDataByUid(context.Background(), &pb.GetUserDataByUidRequest{
		Id: uid,
	})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}

	res.Message = "OK"
	res.Data = response.Data
	return res
}

func (r User) GetRankList() *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	// 获取排行榜
	//reply := cache.QueryRankList()
	//if reply == nil {
	//	res.Code = models.Failed
	//	res.Message = "排行榜获取失败"
	//	logrus.Debugln("排行榜获取失败")
	//	return res
	//}
	//var rankList []models.RankList
	//for i := 0; i < len(reply); i += 2 {
	//	user := models.UserInfo{}
	//	json.Unmarshal([]byte(reply[i]), &user)
	//	item := models.RankList{
	//		Phone:     user.Mobile,
	//		NickName:  user.NickName,
	//		PassCount: user.PassCount,
	//	}
	//	rankList = append(rankList, item)
	//}
	//
	//bytes, err := json.Marshal(rankList)
	//if err != nil {
	//	res.Code = models.Failed
	//	res.Message = err.Error()
	//	logrus.Errorf("排行榜序列化失败:%s", err.Error())
	//	return res
	//}
	//
	//res.Message = "OK"
	//res.Data = string(bytes)
	return res
}

func (r User) GetRecordList(uid int64, page, pageSize int) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}
	// 查询结果
	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewRecordServiceClient(dbConn)
	response, err := client.GetUserRecordList(context.Background(), &pb.GetUserRecordListRequest{
		Uid:      uid,
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Message = "OK"
	res.Data = response
	return res
}

func (r User) GetRecord(id int64) *models.Response {
	return nil
}

// GetUserSolvedList 获取用户解决了哪些题目
func (r User) GetUserSolvedList(uid int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewUserServiceClient(dbConn)
	response, err := client.GetUserSolvedList(context.Background(), &pb.GetUserSolvedListRequest{Uid: uid})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}
	res.Message = "OK"
	res.Data = response.ProblemSolvedList

	return res
}

func (r User) queryUserSolvedListByProblemList(params *models.UPSSParams) ([]int64, error) {
	dbConn, err := registry.NewDBConnection()
	if err != nil {
		logrus.Errorf("db服务连接失败:%s\n", err.Error())
		return nil, err
	}
	defer dbConn.Close()

	client := pb.NewUserServiceClient(dbConn)
	response, err := client.QueryUserSolvedListByProblemIds(context.Background(), &pb.QueryUserSolvedListByProblemIdsRequest{
		Uid:        params.Uid,
		ProblemIds: params.ProblemIds,
	})
	if err != nil {
		return nil, err
	}

	return response.SolvedProblemIds, nil
}
