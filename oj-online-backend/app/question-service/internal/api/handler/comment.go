package handler

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/mq"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

type CommentHandler struct{}

func (receiver CommentHandler) HandleAddComment(form *models.CommentForm) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	pbComment := pb.Comment{
		ObjId:         form.ObjId,
		UserId:        form.UserId,
		UserName:      form.UserName,
		UserAvatarUrl: form.UserAvatarUrl,
		RootId:        form.RootId,
		ReplyId:       form.ReplyId,
		Content:       form.Content,
		Status:        1,
		LikeCount:     0,
		ChildCount:    0,
		Stamp:         form.Stamp,
	}

	// 根评论
	if pbComment.RootId == 0 {
		pbComment.IsRoot = true
		pbComment.ReplyId = 0
	} else {
		// 回复评论
		pbComment.IsRoot = false
		pbComment.ReplyId = form.ReplyId
	}

	if pbComment.Stamp <= 0 {
		pbComment.Stamp = time.Now().Unix()
	}

	// 异步
	option, err := proto.Marshal(&pbComment)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		return res
	}
	productor := mq.Producer{
		Exname:     consts.RabbitMqExchangeName,
		Exkind:     consts.RabbitMqExchangeKind,
		QuName:     consts.RabbitMqCommentQueue,
		RoutingKey: consts.RabbitMqCommentKey,
	}
	if !productor.Publish(option) {
		res.Code = http.StatusInternalServerError
		logrus.Errorln("评论任务提交mq失败")
		return res
	} else {
		res.Code = http.StatusOK
		res.Message = "OK"
		return res
	}
}
