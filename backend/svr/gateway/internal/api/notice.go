package api

import (
	"context"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NoticeHandler struct {
	clients *client.Clients
}

func NewNoticeHandler(clients *client.Clients) *NoticeHandler {
	return &NoticeHandler{
		clients: clients,
	}
}

// 新增公告
func (n *NoticeHandler) HandleAddNotice(ctx *gin.Context) {
	req := &model.AddNoticeRequest{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	uid := ctx.GetInt64("uid")
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := n.clients.Notice.CreateNotice(c, &pb.CreateNoticeRequest{
		Notice: &pb.Notice{
			Title:    req.Title,
			Content:  req.Content,
			CreateBy: uid,
		},
	})
	if err != nil {
		logrus.Errorf("新增公告失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 删除公告
func (n *NoticeHandler) HandleDeleteNotice(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	if err != nil || id <= 0 {
		ResponseBadRequest(ctx, "无效的 id")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = n.clients.Notice.DeleteNotice(c, &pb.DeleteNoticeRequest{Id: id})
	if err != nil {
		logrus.Errorf("删除公告失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 更新公告
func (n *NoticeHandler) HandleUpdateNotice(ctx *gin.Context) {
	req := &model.UpdateNoticeRequest{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err := n.clients.Notice.UpdateNotice(c, &pb.UpdateNoticeRequest{
		Notice: &pb.Notice{
			Id:      req.ID,
			Title:   req.Title,
			Content: req.Content,
		},
	})
	if err != nil {
		logrus.Errorf("更新公告失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 发布公告
func (n *NoticeHandler) HandlePublishNotice(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	if err != nil || id <= 0 {
		ResponseBadRequest(ctx, "无效的 id")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = n.clients.Notice.PublishNotice(c, &pb.PublishNoticeRequest{Id: id})
	if err != nil {
		logrus.Errorf("发布公告失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 隐藏公告
func (n *NoticeHandler) HandleHideNotice(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	if err != nil || id <= 0 {
		ResponseBadRequest(ctx, "无效的 id")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = n.clients.Notice.HideNotice(c, &pb.HideNoticeRequest{Id: id})
	if err != nil {
		logrus.Errorf("隐藏公告失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 获取公告列表
func (n *NoticeHandler) HandleGetNoticeList(ctx *gin.Context) {
	req := &model.GetNoticeListRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := n.clients.Notice.GetNoticeList(c, &pb.GetNoticeListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logrus.Errorf("获取公告列表失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gr)
}

// 获取公告详情
func (n *NoticeHandler) HandleGetNoticeDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		ResponseBadRequest(ctx, "无效的 id")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := n.clients.Notice.GetNotice(c, &pb.GetNoticeRequest{Id: id})
	if err != nil {
		logrus.Errorf("获取公告详情失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gr)
}
