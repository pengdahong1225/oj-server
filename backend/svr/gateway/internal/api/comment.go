package api

import (
	"context"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"time"

	"github.com/gin-gonic/gin"
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
	req := &model.CreateCommentRequest{}
	err := ctx.ShouldBindBodyWithJSON(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = ch.clients.Comment.CreateComment(c, &pb.CreateCommentRequest{
		ObjId:          req.ObjId,
		UserId:         req.UserId,
		UserName:       req.UserName,
		UserAvatarUrl:  req.UserAvatarUrl,
		Content:        req.Content,
		Stamp:          time.Now().Unix(),
		RootId:         req.RootId,
		RootCommentId:  req.RootCommentId,
		ReplyId:        req.ReplyId,
		ReplyCommentId: req.ReplyCommentId,
		ReplyUserName:  req.ReplyUserName,
		Address:        ctx.ClientIP(),
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, nil)
}

// 获取根评论列表
func (ch *CommentHandler) HandleGetRootCommentList(ctx *gin.Context) {
	req := &model.QueryRootCommentListRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := ch.clients.Comment.QueryRootComment(c, &pb.QueryRootCommentRequest{
		ObjId:    req.ObjId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 获取子评论列表
func (ch *CommentHandler) HandleGetChildCommentList(ctx *gin.Context) {
	req := &model.QueryChildCommentListRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := ch.clients.Comment.QueryChildComment(c, &pb.QueryChildCommentRequest{
		ObjId:          req.ObjId,
		RootId:         req.RootId,
		RootCommentId:  req.RootCommentId,
		ReplyId:        req.ReplyId,
		ReplyCommentId: req.ReplyCommentId,
		Cursor:         req.Cursor,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 点赞评论
func (ch *CommentHandler) HandleLikeComment(ctx *gin.Context) {
	req := &model.CommentLikeRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = ch.clients.Comment.CommentLike(c, &pb.CommentLikeRequest{
		ObjId:     req.ObjId,
		CommentId: req.CommentId,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, nil)
}

// 删除评论
func (ch *CommentHandler) HandleDeleteComment(ctx *gin.Context) {
	req := &model.DeleteCommentRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = ch.clients.Comment.DeleteComment(c, &pb.DeleteCommentRequest{
		ObjId:     req.ObjId,
		CommentId: req.CommentId,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, nil)
}
