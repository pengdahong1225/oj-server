package logic

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"question-service/global"
	"question-service/middlewares"
	"question-service/models"
	pb "question-service/proto"
	"strconv"
	"time"
)

// 签名过期时间
var (
	jwtTimeOut int64 = 60 * 60 * 24 * 7
	issuer           = "Messi"
)

func RegistryHandler(ctx *gin.Context, form *models.RegistryForm) {
	// 验证码校验
	redisConn := global.RedisPoolInstance.Get()
	defer redisConn.Close()
	if c, err := redis.String(redisConn.Do("Get", form.Phone)); err != nil {
		logrus.Debugln(err.Error())
		ctx.JSON(http.StatusNoContent, gin.H{
			"msg": "验证码不存在",
		})
	} else if c != form.SmsCode {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
	}

	// 注册
	dbConn, err := global.NewDBConnection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
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
	dbConn, err := global.NewDBConnection()
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

func GetUserDetail(ctx *gin.Context) {
	phone, _ := strconv.ParseInt(ctx.Query("phone"), 10, 64)
	dbConn, err := global.NewDBConnection()
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
	conn := global.RedisPoolInstance.Get()
	defer conn.Close()

	reply, err := redis.Strings(conn.Do("zrange", "rank", 0, -1))
	if err != nil {
		logrus.Debugln(err.Error())
		ctx.JSON(http.StatusNoContent, gin.H{
			"msg": "排行榜获取失败",
		})
		return
	}

	ranklist := make([]models.RankList, 20)
	for i := 0; i < len(reply); i += 2 {
		user := models.UserInfo{}
		json.Unmarshal([]byte(reply[i]), &user)
		item := models.RankList{
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

func GetSubmitRecord(ctx *gin.Context) {
	// 获取提交记录
}
