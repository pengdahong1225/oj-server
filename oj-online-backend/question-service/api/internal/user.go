package internal

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"question-service/api/proto"
	"question-service/middlewares"
	"question-service/models"
	"question-service/services/redis"
	"question-service/services/registry"
	"question-service/services/sms"
	"question-service/settings"
	"question-service/utils"
	"question-service/views"
	"strconv"
	"time"
)

// 签名过期时间
var (
	jwtTimeOut int64 = 60 * 60 * 24 * 7
	issuer           = "Messi"
)

func RegistryHandler(ctx *gin.Context, form *models.RegistryForm) {
	var cms_phone = "18048155008"
	// 验证码校验
	redisConn := redis.NewConn()
	defer redisConn.Close()
	if c, err := redigo.String(redisConn.Do("Get", cms_phone)); err != nil {
		logrus.Debugln(err.Error())
		ctx.JSON(http.StatusNoContent, gin.H{
			"msg": "验证码不存在",
		})
		return
	} else if c != form.SmsCode {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 注册
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
		return
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	phone, _ := strconv.ParseInt(form.Phone, 10, 64)
	data := &pb.UserInfo{
		Phone:    phone,
		Password: form.PassWord,
		Nickname: form.NickName,
		Email:    form.Email,
		Gender:   int32(form.Gender),
		Role:     int32(form.Role),
		HeadUrl:  form.HeadUrl,
	}
	request := &pb.CreateUserRequest{Data: data}

	response, err := client.CreateUserData(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
		"id":  response.Id,
	})
}

func LoginHandler(ctx *gin.Context, form *models.LoginFrom) {
	// 注册
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)

	phone, _ := strconv.ParseInt(form.Phone, 10, 64)
	request := &pb.GetUserRequest{
		Phone: phone,
	}
	response, err := client.GetUserData(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	// 校验密码
	if response.Data.Password != form.PassWord {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "密码错误",
		})
		return
	}

	// 生成token
	j := middlewares.NewJWT()
	// 设置 payload有效载荷
	claims := &middlewares.UserClaims{
		Phone:     response.Data.Phone,
		Authority: response.Data.Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              // 签名生效时间
			ExpiresAt: time.Now().Unix() + jwtTimeOut, // 7天过期
			Issuer:    issuer,                         // 签名机构
		},
	}
	xtoken, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	// 序列化
	data, _ := protojson.Marshal(response.Data)
	ctx.JSON(http.StatusOK, gin.H{
		"data":   data,
		"xtoken": xtoken,
	})
}

func GetUserDetail(ctx *gin.Context, phone int64) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetUserData(context.Background(), &pb.GetUserRequest{
		Phone: phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	data, _ := protojson.Marshal(response.Data)
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func GetRankList(ctx *gin.Context) {
	// 获取排行榜
	conn := redis.NewConn()
	defer conn.Close()

	reply, err := redigo.Strings(conn.Do("zrange", "rank", 0, -1))
	if err != nil {
		logrus.Debugln(err.Error())
		ctx.JSON(http.StatusNoContent, gin.H{
			"msg": "排行榜获取失败",
		})
		return
	}

	ranklist := make([]views.RankList, 20)
	for i := 0; i < len(reply); i += 2 {
		user := views.UserInfo{}
		json.Unmarshal([]byte(reply[i]), &user)
		item := views.RankList{
			Phone:     user.Phone,
			NickName:  user.NickName,
			PassCount: user.PassCount,
		}
		ranklist = append(ranklist, item)
	}

	// 序列化
	data, _ := json.Marshal(ranklist)
	ctx.JSON(http.StatusOK, gin.H{
		"rankList": data,
	})
}

func GetSubmitRecord(ctx *gin.Context, userId int64) {
	// 获取提交记录
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()
	client := pb.NewDBServiceClient(dbConn)

	response, err := client.GetUserSubmitRecord(context.Background(), &pb.GetUserSubmitRecordRequest{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 序列化
	data, _ := json.Marshal(response.Data)
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func SendSmsCode(phone string) error {
	// 生成随机数
	c := utils.GenerateSmsCode(6)
	expire := 180 // 3min过期
	param := map[string]string{
		"code": c,
	}
	data, _ := json.Marshal(param)

	// 调用第三方服务发送
	if err := sms.Send(data, phone); err != nil {
		return err
	}

	// 缓存验证码
	redisConn := redis.NewConn()
	defer redisConn.Close()
	if _, err := redisConn.Do("Set", phone, c, "ex", expire); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
