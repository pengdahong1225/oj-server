package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/global"
	"oj-server/pkg/proto/pb"
	"oj-server/pkg/registry"
	"oj-server/svr/gateway/internal/model"
	"time"
)

func HandleCreateComment(ctx *gin.Context) {
	// 参数校验
	form, ok := validateWithJson(ctx, model.CreateCommentForm{})
	if !ok {
		return
	}

	// 调用评论服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.CreateCommentRequest{
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
	}
	_, err = client.CreateComment(ctx, request)
	if err != nil {
		logrus.Errorf("创建评论失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "创建评论失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", nil)
}
func HandleGetRootCommentList(ctx *gin.Context) {
	// 参数校验
	var params model.QueryRootCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "参数校验失败", nil)
		return
	}

	// 调用comment服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.QueryRootCommentRequest{
		ObjId:    params.ObjId,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	resp, err := client.QueryRootComment(context.Background(), request)
	if err != nil {
		logrus.Errorf("获取root评论列表失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取root评论列表失败", nil)
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

	ResponseWithJson(ctx, http.StatusOK, "success", result)
}
func HandleGetChildCommentList(ctx *gin.Context) {
	// 参数校验
	var params model.QueryChildCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "参数校验失败", nil)
		return
	}

	// 调用comment服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.QueryChildCommentRequest{
		ObjId:          params.ObjId,
		RootId:         params.RootId,
		RootCommentId:  params.RootCommentId,
		ReplyId:        params.ReplyId,
		ReplyCommentId: params.ReplyCommentId,
		Cursor:         params.Cursor,
	}
	resp, err := client.QueryChildComment(ctx, request)
	if err != nil {
		logrus.Errorf("获取child评论列表失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取child评论列表失败", nil)
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

	ResponseWithJson(ctx, http.StatusOK, "success", result)
}
func HandleLikeComment(ctx *gin.Context) {
	form, ok := validateWithForm(ctx, model.CommentLikeForm{})
	if !ok {
		return
	}
	// 调用comment服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.CommentLikeRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	}
	if _, err = client.CommentLike(ctx, request); err != nil {
		logrus.Errorf("评论点赞: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "评论点赞失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", nil)
}
func HandleDeleteComment(ctx *gin.Context) {
	form, ok := validateWithForm[model.DeleteCommentForm](ctx, model.DeleteCommentForm{})
	if !ok {
		return
	}
	// 调用comment服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.DeleteCommentRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	}
	if _, err = client.DeleteComment(ctx, request); err != nil {
		logrus.Errorf("删除评论失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "删除评论失败", nil)
		return
	}

	ResponseWithJson(ctx, http.StatusOK, "success", nil)
}
