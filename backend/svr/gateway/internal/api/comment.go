package api

import (
	"context"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CommentHandler struct {
	clients *client.Clients
}

func NewCommentHandler(clients *client.Clients) *CommentHandler {
	return &CommentHandler{
		clients: clients,
	}
}

// 创建评论
func (ch *CommentHandler) HandleCreateComment(ctx *gin.Context) {
	form, ok := validateWithJson(ctx, model.CreateCommentForm{})
	if !ok {
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := ch.clients.Comment.CreateComment(c, &pb.CreateCommentRequest{
		ObjId:          form.ObjId,
		UserId:         form.UserId,
		UserName:       form.UserName,
		UserAvatarUrl:  form.UserAvatarUrl,
		Content:        form.Content,
		Stamp:          time.Now().Unix(),
		RootId:         form.RootId,
		RootCommentId:  form.RootCommentId,
		ReplyId:        form.ReplyId,
		ReplyCommentId: form.ReplyCommentId,
		ReplyUserName:  form.ReplyUserName,
		Address:        ctx.ClientIP(),
	})
	if err != nil {
		logrus.Errorf("创建评论失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 获取根评论列表
func (ch *CommentHandler) HandleGetRootCommentList(ctx *gin.Context) {
	var params model.QueryRootCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseBadRequest(ctx, "参数校验失败")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := ch.clients.Comment.QueryRootComment(c, &pb.QueryRootCommentRequest{
		ObjId:    params.ObjId,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
	if err != nil {
		logrus.Errorf("获取root评论列表失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	result := &model.QueryRootCommentListResult{
		Total: int64(resp.Total),
	}
	result.List = make([]*model.Comment, 0, len(resp.Data))
	for _, pbComment := range resp.Data {
		comment := &model.Comment{}
		comment.FromPbComment(pbComment)
		result.List = append(result.List, comment)
	}
	ResponseOK(ctx, result)
}

// 获取子评论列表
func (ch *CommentHandler) HandleGetChildCommentList(ctx *gin.Context) {
	var params model.QueryChildCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseBadRequest(ctx, "参数校验失败")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := ch.clients.Comment.QueryChildComment(c, &pb.QueryChildCommentRequest{
		ObjId:          params.ObjId,
		RootId:         params.RootId,
		RootCommentId:  params.RootCommentId,
		ReplyId:        params.ReplyId,
		ReplyCommentId: params.ReplyCommentId,
		Cursor:         params.Cursor,
	})
	if err != nil {
		logrus.Errorf("获取child评论列表失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	result := model.QueryChildCommentListResult{
		Total:  int64(resp.Total),
		Cursor: resp.Cursor,
	}
	result.List = make([]*model.Comment, 0, len(resp.Data))
	for _, pbComment := range resp.Data {
		comment := &model.Comment{}
		comment.FromPbComment(pbComment)
		result.List = append(result.List, comment)
	}

	ResponseOK(ctx, result)
}

// 点赞评论
func (ch *CommentHandler) HandleLikeComment(ctx *gin.Context) {
	form, ok := validateWithForm(ctx, model.CommentLikeForm{})
	if !ok {
		return
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := ch.clients.Comment.CommentLike(c, &pb.CommentLikeRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	})
	if err != nil {
		logrus.Errorf("评论点赞: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 删除评论
func (ch *CommentHandler) HandleDeleteComment(ctx *gin.Context) {
	form, ok := validateWithForm[model.DeleteCommentForm](ctx, model.DeleteCommentForm{})
	if !ok {
		return
	}
	// 调用comment服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := ch.clients.Comment.DeleteComment(c, &pb.DeleteCommentRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	})
	if err != nil {
		logrus.Errorf("删除评论失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, nil)
}
