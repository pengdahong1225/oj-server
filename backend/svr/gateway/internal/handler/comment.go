package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/global"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/model"
	"oj-server/svr/gateway/internal/service"
	"time"
)

func HandleCreateComment(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}
	// 参数校验
	form, ok := validateWithJson(ctx, model.CreateCommentForm{})
	if !ok {
		return
	}
	pbComment := &pb.Comment{
		ObjId:         form.ObjId,
		UserId:        form.UserId,
		UserName:      form.UserName,
		UserAvatarUrl: form.UserAvatarUrl,
		Content:       form.Content,
		Status:        1,
		PubStamp:      time.Now().Unix(),
		PubRegion:     "",
	}
	// 非楼主评论
	if form.RootId > 0 && form.RootCommentId > 0 {
		pbComment.IsRoot = 0
		pbComment.RootId = form.RootId
		pbComment.RootCommentId = form.RootCommentId
		if form.ReplyId > 0 && form.ReplyCommentId > 0 {
			pbComment.ReplyId = form.ReplyId
			pbComment.ReplyCommentId = form.ReplyCommentId
			pbComment.ReplyUserName = form.ReplyUserName
		} else {
			// 默认回复楼主
			pbComment.ReplyId = form.RootId
			pbComment.ReplyCommentId = form.RootCommentId
			pbComment.ReplyUserName = form.ReplyUserName
		}
	} else {
		// 楼主评论
		pbComment.RootId = 0
		pbComment.RootCommentId = 0
		pbComment.IsRoot = 1
		pbComment.ReplyId = 0
		pbComment.ReplyCommentId = 0
	}
	// 查询ip归属地
	ip := ctx.ClientIP()
	info, err := service.QueryIpGeolocation(ip)
	if err != nil {
		logrus.Errorf("查询ip归属地失败,ip:%s,err:%s", ip, err.Error())
		info = &model.IPInfoResp{
			RegionName: "未知地区",
		}
	}
	pbComment.PubRegion = info.RegionName

	// 发送到rabbitmq
	if err = service.SendComment2MQ(pbComment); err != nil {
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
	}
	ctx.JSON(http.StatusOK, resp)
}
func HandleGetRootCommentList(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}
	// 参数校验
	var params model.QueryRootCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "参数校验失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用comment服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败: %s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.QueryRootCommentRequest{
		ObjId:    params.ObjId,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	rpc_resp, err := client.QueryRootComment(context.Background(), request)
	if err != nil {
		logrus.Errorf("获取root评论列表失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取评论列表失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := &model.QueryRootCommentListResult{
		Total: int64(rpc_resp.Total),
	}
	data.List = make([]*model.Comment, 0, len(rpc_resp.Data))
	for _, pbComment := range rpc_resp.Data {
		comment := &model.Comment{}
		comment.FromPbComment(pbComment)
		data.List = append(data.List, comment)
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}
func HandleGetChildCommentList(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}
	// 参数校验
	var params model.QueryChildCommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "参数校验失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用comment服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败: %s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
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
	rpc_resp, err := client.QueryChildComment(ctx, request)
	if err != nil {
		logrus.Errorf("获取child评论列表失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取child评论列表失败"
		ctx.JSON(http.StatusInternalServerError, resp)
	}
	data := model.QueryChildCommentListResult{
		Total:  int64(rpc_resp.Total),
		Cursor: rpc_resp.Cursor,
	}
	data.List = make([]*model.Comment, 0, len(rpc_resp.Data))
	for _, pbComment := range rpc_resp.Data {
		comment := &model.Comment{}
		comment.FromPbComment(pbComment)
		data.List = append(data.List, comment)
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}
func HandleLikeComment(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}

	form, ok := validateWithForm(ctx, model.CommentLikeForm{})
	if !ok {
		return
	}
	// 调用comment服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败: %s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.CommentLikeRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	}
	if _, err = client.CommentLike(ctx, request); err != nil {
		logrus.Errorf("评论点赞: %s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "评论点赞失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleDeleteComment(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}

	form, ok := validateWithForm[model.DeleteCommentForm](ctx, model.DeleteCommentForm{})
	if !ok {
		return
	}
	// 调用comment服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败: %s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewCommentServiceClient(conn)
	request := &pb.DeleteCommentRequest{
		ObjId:     form.ObjId,
		CommentId: form.CommentId,
	}
	if _, err = client.DeleteComment(ctx, request); err != nil {
		logrus.Errorf("删除评论失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "删除失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
