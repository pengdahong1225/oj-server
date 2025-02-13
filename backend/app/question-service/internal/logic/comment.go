package logic

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

type CommentLogic struct{}

func (receiver CommentLogic) OnAddComment(form *models.AddCommentForm, region string) *models.Response {
	res := &models.Response{
		Code:    models.Success,
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
	}
	pbComment.PubRegion = region

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

	if pbComment.PubStamp <= 0 {
		pbComment.PubStamp = time.Now().Unix()
	}

	// 异步
	msg, err := proto.Marshal(&pbComment)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}
	productor := mq.NewProducer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqCommentQueue,
		consts.RabbitMqCommentKey,
	)
	if !productor.Publish(msg) {
		res.Code = models.Failed
		res.Message = "任务提交mq失败"
		logrus.Errorln("评论任务提交mq失败")
		return res
	} else {
		res.Message = "OK"
		return res
	}
}

func (receiver CommentLogic) OnQueryRootCommentList(params *models.RootCommentListQueryParams) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	request := &pb.QueryRootCommentRequest{
		ObjId:    params.ObjId,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewCommentServiceClient(dbConn)
	response, err := client.QueryRootComment(context.Background(), request)
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

func (receiver CommentLogic) OnQueryChildCommentList(params *models.ChildCommentListQueryParams) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	request := &pb.QueryChildCommentRequest{
		ObjId:          params.ObjId,
		RootId:         params.RootId,
		RootCommentId:  params.RootCommentId,
		ReplyId:        params.ReplyId,
		ReplyCommentId: params.ReplyCommentId,
		Cursor:         params.Cursor,
	}
	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewCommentServiceClient(dbConn)

	response, err := client.QueryChildComment(context.Background(), request)
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

func (receiver CommentLogic) OnCommentLike(form *models.CommentLikeForm) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	request := &pb.CommentLikeRequest{
		CommentId: form.CommentId,
		ObjId:     form.ObjId,
	}
	dbConn, err := registry.NewDBConnection()
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewCommentServiceClient(dbConn)

	_, err = client.CommentLike(context.Background(), request)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Message = "OK"
	return res
}
