package logic

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/global"
	"question-service/models"
	pb "question-service/proto"
)

func GetUserList(ctx *gin.Context, cursor int32) {
	dbConn, err := global.NewDBConnection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.GetUserListRequest{Cursor: cursor}
	response, err := client.GetUserList(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var userList []*models.UserInfo
	for _, u := range response.Data {
		userList = append(userList, &models.UserInfo{
			Phone:       u.Phone,
			NickName:    u.Nickname,
			Email:       u.Email,
			Gender:      u.Gender,
			Role:        u.Role,
			HeadUrl:     u.HeadUrl,
			PassCount:   u.PassCount,
			CreateAt:    u.CreateAt,
			SubmitCount: u.SubmitCount,
		})
	}
	data, _ := json.Marshal(userList)
	ctx.JSON(http.StatusOK, gin.H{
		"data":   data,
		"cursor": response.Cursor,
	})
}
