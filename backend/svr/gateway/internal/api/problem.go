package api

import (
	"context"
	"fmt"
	"io"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"path/filepath"
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
		Title:       req.Title,
		Description: req.Description,
		Level:       req.Level,
		Tags:        req.Tags,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}
	ResponseOK(ctx, gr)
}

// 处理题目配置文件上传
func (p *ProblemHandler) HandleUploadConfig(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil || problemId <= 0 {
		ResponseBadRequest(ctx, "无效的 problem_id")
		return
	}

	// 获取文件
	fileHeader, err := ctx.FormFile("config_file")
	if err != nil {
		ResponseBadRequest(ctx, "需要config_file字段")
		return
	}
	if filepath.Ext(fileHeader.Filename) != ".json" {
		ResponseBadRequest(ctx, "文件格式错误")
		return
	}
	// 限制文件大小（默认 8MB，可调整）
	maxSize := int64(8 << 20) // 8MB
	if fileHeader.Size > maxSize {
		ResponseBadRequest(ctx, "文件过大")
		return
	}
	// 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		ResponseError(ctx, pb.Error_EN_Failed, "文件打开失败")
		return
	}

	// 创建gRPC流
	md := metadata.Pairs("problem_id", fmt.Sprintf("%d", problemId))
	c, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx.Request.Context(), md), consts.GrpcStreamTimeout)
	defer cancel()
	stream, err := p.clients.Problem.UploadConfig(c)
	if err != nil {
		logrus.Errorf("文件上传失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	// 分片读取并流式传输
	buffer := make([]byte, 1<<20) // 1MB分片
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("file read error:%s", err.Error())
			_ = stream.CloseSend()
			ResponseError(ctx, pb.Error_EN_Failed, "file read failed")
			return
		}
		// 发送分片到Problem服务
		if err = stream.Send(&pb.UploadConfigFileChunk{
			Content:   buffer[:n],
			ProblemId: problemId,
		}); err != nil {
			logrus.Errorf("file send error:%s", err.Error())
			_ = stream.CloseSend()
			ResponseWithGrpcError(ctx, err)
			return
		}
	}
	// 获取Problem服务响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		logrus.Errorf("file save error:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	if resp.Size != fileHeader.Size {
		logrus.Errorf("file save error, file size[%d] != save size[%d]", fileHeader.Size, resp.Size)
		ResponseError(ctx, pb.Error_EN_Failed, "file save failed")
		return
	}

	ResponseOK(ctx, nil)
}

// 下载题目配置文件
func (p *ProblemHandler) HandleDownloadConfig(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil || problemId <= 0 {
		ResponseBadRequest(ctx, "无效的 problem_id")
		return
	}
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
			Id:          req.ProblemID,
			Title:       req.Title,
			Description: req.Desc,
			Level:       req.Level,
			Tags:        req.Tags,
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
