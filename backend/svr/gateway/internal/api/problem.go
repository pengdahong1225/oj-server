package api

import (
	"context"
	"fmt"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type ProblemHandler struct {
	clients *client.Clients
}

func NewProblemHandler(clients *client.Clients) *ProblemHandler {
	return &ProblemHandler{
		clients: clients,
	}
}

// 获取题目标签列表
func (p *ProblemHandler) HandleGetProblemTagList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Problem.GetTagList(c, nil)
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 获取题目列表
func (p *ProblemHandler) HandleGetProblemList(ctx *gin.Context) {
	req := &model.QueryProblemListRequest{}
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	role, exist := ctx.Get("role")
	if !exist {
		ResponseUnauthorized(ctx, "invalid role")
		return
	}

	// 调用problem服务
	md := metadata.Pairs("role", fmt.Sprintf("%d", role.(int32)))
	c, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx.Request.Context(), md), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Problem.GetProblemList(c, &pb.GetProblemListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Tag:      req.Tag,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 获取题目详情
func (p *ProblemHandler) HandleGetProblemDetail(ctx *gin.Context) {
	problem_id, err := strconv.ParseInt(ctx.Query("problem_id"), 10, 64)
	if err != nil {
		ResponseBadRequest(ctx, "problem_id validate err")
		return
	}
	// 调用problem服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Problem.GetProblemDetail(c, &pb.GetProblemRequest{
		Id: problem_id,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gr)
}

// 创建题目信息
func (p *ProblemHandler) HandleCreateProblem(ctx *gin.Context) {
	req := &model.CreateProblemRequest{}
	err := ctx.ShouldBindBodyWithJSON(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	// 调用problem服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Problem.CreateProblem(c, &pb.CreateProblemRequest{
		Title: req.Title,
		Level: req.Level,
		Tags:  req.Tags,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gr)
}

// 处理发布题目
func (p *ProblemHandler) HandlePublishProblem(ctx *gin.Context) {
	problem_id, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil {
		ResponseBadRequest(ctx, "无效的 problem_id")
		return
	}
	// 调用problem服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = p.clients.Problem.PublishProblem(c, &pb.PublishProblemRequest{
		Id: problem_id,
	})
	if err != nil {
		logrus.Errorf("发布失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 处理隐藏题目
func (p *ProblemHandler) HandleHideProblem(ctx *gin.Context) {
	problem_id, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil {
		ResponseBadRequest(ctx, "无效的 problem_id")
		return
	}
	// 调用problem服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = p.clients.Problem.HideProblem(c, &pb.HideProblemRequest{
		Id: problem_id,
	})
	if err != nil {
		logrus.Errorf("隐藏失败: %s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 处理删除题目
func (p *ProblemHandler) HandleDeleteProblem(ctx *gin.Context) {
	problem_id, err := strconv.ParseInt(ctx.Query("problem_id"), 10, 64)
	if err != nil {
		ResponseBadRequest(ctx, "无效的 problem_id")
		return
	}
	// 调用problem服务
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = p.clients.Problem.DeleteProblem(c, &pb.DeleteProblemRequest{
		Id: problem_id,
	})
	if err != nil {
		logrus.Errorf("删除题目失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 处理更新题目信息
func (p *ProblemHandler) HandleUpdateProblem(ctx *gin.Context) {
	req := &model.UpdateProblemRequest{}
	err := ctx.ShouldBindBodyWithJSON(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	_, err = p.clients.Problem.UpdateProblem(c, &pb.UpdateProblemRequest{
		Data: &pb.Problem{
			Id:    req.ProblemID,
			Title: req.Title,
			Level: req.Level,
			Tags:  req.Tags,
		},
	})
	if err != nil {
		logrus.Errorf("更新题目失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, nil)
}

// 提交代码
func (p *ProblemHandler) HandleSubmitProblem(ctx *gin.Context) {
	req := &model.SubmitRequest{}
	err := ctx.ShouldBindBodyWithJSON(req)
	if err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	uid := ctx.GetInt64("uid")
	md := metadata.Pairs("uid", fmt.Sprintf("%d", uid))
	c, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx.Request.Context(), md), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Problem.SubmitProblem(c, &pb.SubmitProblemRequest{
		ProblemId: req.ProblemID,
		Title:     req.Title,
		Lang:      req.Lang,
		Code:      req.Code,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 获取判题任务结果
func (p *ProblemHandler) HandleGetSubmitResult(ctx *gin.Context) {
	taskId, exists := ctx.GetQuery("task_id")
	if !exists || taskId == "" {
		ResponseBadRequest(ctx, "task_id invalid")
		return
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := p.clients.Record.QueryJudgeResult(c, &pb.QueryJudgeResultRequest{
		TaskId: taskId,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 处理获取用户的AC题目列表
func (p *ProblemHandler) HandleGetUserSolvedList(ctx *gin.Context) {}
