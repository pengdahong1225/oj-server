package logic

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/svc/mq"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

type CommentLogic struct{}

func (receiver CommentLogic) HandleAddComment(form *models.CommentForm) *models.Response {
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
		Content:       form.Content,
		Status:        1,
		LikeCount:     0,
		ChildCount:    0,
		Stamp:         form.Stamp,
	}

	// 非楼主评论
	if form.RootId > 0 && form.RootCommentId > 0 {
		pbComment.RootId = form.RootId
		pbComment.RootCommentId = form.RootCommentId
		pbComment.IsRoot = 0
		if form.ReplyId > 0 && form.ReplyCommentId > 0 {
			pbComment.ReplyId = form.ReplyId
			pbComment.ReplyCommentId = form.ReplyCommentId
		} else {
			// 默认回复楼主
			pbComment.ReplyId = form.RootId
			pbComment.ReplyCommentId = form.RootCommentId
		}
	} else {
		// 楼主评论
		pbComment.RootId = 0
		pbComment.RootCommentId = 0
		pbComment.IsRoot = 1
		pbComment.ReplyId = 0
		pbComment.ReplyCommentId = 0
	}

	if pbComment.Stamp <= 0 {
		pbComment.Stamp = time.Now().Unix()
	}

	// 异步
	msg, err := proto.Marshal(&pbComment)
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
	if !productor.Publish(msg) {
		res.Code = http.StatusInternalServerError
		logrus.Errorln("评论任务提交mq失败")
		return res
	} else {
		res.Code = http.StatusOK
		res.Message = "OK"
		return res
	}
}
