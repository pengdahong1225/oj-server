package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"oj-server/global"
	"oj-server/pkg/proto/pb"
	"oj-server/pkg/registry"
	"oj-server/svr/gateway/internal/model"
	"path/filepath"
	"strconv"
)

// 获取题目标签列表
func HandleGetProblemTagList(ctx *gin.Context) {
	// 调用problem服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	resp, err := client.GetTagList(context.Background(), &empty.Empty{})
	if err != nil {
		logrus.Errorf("获取题目标签列表失败: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取题目标签列表失败", nil)
		return
	}

	ResponseWithJson(ctx, http.StatusOK, "success", resp.Data)
}

// 获取题目列表
func HandleGetProblemList(ctx *gin.Context) {
	// 查询参数校验
	var params model.QueryProblemListParams
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		logrus.Errorf("参数校验失败: %s", err.Error())
		ResponseWithJson(ctx, http.StatusBadRequest, "参数校验失败", nil)
		return
	}

	// 调用problem服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.GetProblemListRequest{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Tag:      params.Tag,
	}
	resp, err := client.GetProblemList(context.Background(), req)
	if err != nil {
		logrus.Errorf("problem服务获取题目列表失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取题目列表失败", nil)
		return
	}

	result := &model.QueryProblemListResult{
		Total: resp.Total,
	}
	result.List = make([]*model.Problem, len(resp.Data))
	for i, v := range resp.Data {
		result.List[i] = &model.Problem{
			ID:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Level:       v.Level,
			Tags:        v.Tags,
			Status:      v.Status,
			CreateAt:    v.CreateAt,
		}
	}

	ResponseWithJson(ctx, http.StatusOK, "success", result)
}

// 获取题目详情
func HandleGetProblemDetail(ctx *gin.Context) {
	// 查询参数
	problem_id, err := strconv.ParseInt(ctx.Query("problem_id"), 10, 64)
	if err != nil {
		logrus.Errorf("problem_id validate err: %s", err.Error())
		ResponseWithJson(ctx, http.StatusBadRequest, "problem_id validate err", nil)
		return
	}
	// 调用problem服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	resp, err := client.GetProblemDetail(context.Background(), &pb.GetProblemRequest{
		Id: problem_id,
	})
	if err != nil {
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取题目详情失败", nil)
		return
	}
	problem := &model.Problem{
		ID:          resp.Problem.Id,
		Title:       resp.Problem.Title,
		Description: resp.Problem.Description,
		Level:       resp.Problem.Level,
		Tags:        resp.Problem.Tags,
		Status:      resp.Problem.Status,
		CreateAt:    resp.Problem.CreateAt,
	}
	ResponseWithJson(ctx, http.StatusOK, "success", problem)
}

// 提交代码
func HandleSubmitProblem(ctx *gin.Context) {
	// 表单验证
	form, ok := validateWithJson(ctx, model.SubmitForm{})
	if !ok {
		return
	}

	// 获取元数据
	uid := ctx.GetInt64("uid")

	// 调用题目服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.SubmitProblemRequest{
		ProblemId: form.ProblemID,
		Title:     form.Title,
		Lang:      form.Lang,
		Code:      form.Code,
	}
	resp, err := client.SubmitProblem(context.WithValue(context.Background(), "uid", uid), req)
	if err != nil {
		logrus.Errorf("problem服务提交代码失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "提交失败", nil)
		return
	}

	result := &model.SubmitResult{
		TaskId: resp.TaskId,
	}
	ResponseWithJson(ctx, http.StatusOK, "success", result)
}

// 创建题目信息
func HandleCreateProblem(ctx *gin.Context) {
	form, ret := validateWithJson(ctx, model.CreateProblemForm{})
	if !ret {
		return
	}
	// 调用problem服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	request := &pb.CreateProblemRequest{
		Title:       form.Title,
		Description: form.Description,
		Level:       form.Level,
		Tags:        form.Tags,
	}
	resp, err := client.CreateProblem(context.Background(), request)
	if err != nil {
		logrus.Errorf("创建题目失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "创建题目失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", gin.H{
		"problem_id": resp.Id,
	})
}

// 处理题目配置文件上传
func HandleUploadConfig(ctx *gin.Context) {
	// 获取元数据
	problemId, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil || problemId <= 0 {
		ResponseWithJson(ctx, http.StatusBadRequest, "无效的 problem_id", nil)
		return
	}

	// 获取文件
	fileHeader, err := ctx.FormFile("config_file")
	if err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "需要config_file字段", nil)
		return
	}
	if filepath.Ext(fileHeader.Filename) != ".json" {
		ResponseWithJson(ctx, http.StatusBadRequest, "文件格式错误", nil)
		return
	}
	// 限制文件大小（默认 8MB，可调整）
	maxSize := int64(8 << 20) // 8MB
	if fileHeader.Size > maxSize {
		ResponseWithJson(ctx, http.StatusBadRequest, "文件大小超过限制", nil)
		return
	}
	// 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "文件打开失败", nil)
		return
	}
	// 创建gRPC流
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	stream, err := client.UploadConfig(ctx)
	if err != nil {
		logrus.Errorf("文件上传失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "文件上传失败", nil)
		return
	}
	// 分片读取并流式传输
	buffer := make([]byte, 0, 1<<20) // 1MB分片
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("file read error:%s", err.Error())
			_ = stream.CloseSend()
			ResponseWithJson(ctx, http.StatusInternalServerError, "文件上传失败", nil)
			return
		}
		// 发送分片到Problem服务
		if err = stream.Send(&pb.UploadConfigFileChunk{
			Content:   buffer[:n],
			ProblemId: problemId,
		}); err != nil {
			logrus.Errorf("file send error:%s", err.Error())
			_ = stream.CloseSend()
			ResponseWithJson(ctx, http.StatusInternalServerError, "文件上传失败", nil)
			return
		}
	}
	// 获取Problem服务响应
	_, err = stream.CloseAndRecv()
	if err != nil {
		logrus.Errorf("file save error:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "文件保存失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", nil)
}

// 处理发布题目
func HandlePublishProblem(ctx *gin.Context) {
	// 获取元数据
	problem_id, err := strconv.ParseInt(ctx.PostForm("problem_id"), 10, 64)
	if err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "无效的 problem_id", nil)
		return
	}
	// 调用problem服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	if _, err = client.PublishProblem(ctx, &pb.PublishProblemRequest{
		Id: problem_id,
	}); err != nil {
		logrus.Errorf("发布失败: %s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "发布失败", nil)
		return
	}
	ResponseWithJson(ctx, http.StatusOK, "success", nil)
}

// 处理删除题目
func HandleDeleteProblem(ctx *gin.Context) {}

// 处理更新题目信息
func HandleUpdateProblem(ctx *gin.Context) {}
