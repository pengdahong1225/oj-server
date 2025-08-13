package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"oj-server/global"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/define"
	"path/filepath"
	"strconv"
)

// 处理获取标签列表
func HandleGetTagList(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 调用problem服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	rpc_resp, err := client.GetTagList(context.Background(), &empty.Empty{})
	if err != nil {
		
	}
}

// 处理获取题目列表
func HandleGetProblemList(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
	tag := ctx.Query("tag")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "页码参数错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "页大小参数错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 调用problem服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.GetProblemListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  keyWord,
		Tag:      tag,
	}
	rps_resp, err := client.GetProblemList(context.Background(), req)
	if err != nil {
		logrus.Errorf("problem服务获取题目列表失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取题目列表失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = rps_resp
	resp.Message = "获取题目列表成功"
	ctx.JSON(http.StatusOK, resp)
	return
}

// 处理获取题目详情
func HandleGetProblemDetail(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 查询参数
	idStr := ctx.Query("problem_id")
	if idStr == "" {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "题目id不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	// 调用problem服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.GetProblemRequest{
		Id: id,
	}
	rpc_resp, err := client.GetProblemData(context.Background(), req)
	if err != nil {
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取题目详情失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.ErrCode = pb.Error_EN_Success
	resp.Data = rpc_resp.Data
	resp.Message = "获取题目详情成功"
	ctx.JSON(http.StatusOK, resp)
}

// 处理提交代码
// 判断“用户”是否处于判题状态？true就拒绝
// 用户提交了题目就立刻返回，并给题目设置状态
// 客户端通过其他接口轮询题目结果
func HandleSubmitProblem(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, define.SubmitForm{})
	if !ret {
		return
	}
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 获取元数据
	uid := ctx.GetInt64("uid")

	// 调用题目服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.SubmitProblemRequest{
		ProblemId: form.ProblemID,
		Title:     form.Title,
		Lang:      form.Lang,
		Code:      form.Code,
	}
	rpc_resp, err := client.SubmitProblem(context.WithValue(context.Background(), "uid", uid), req)
	if err != nil {
		logrus.Errorf("problem服务提交代码失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "提交代码失败"
		ctx.JSON(http.StatusOK, resp)
	}
	resp.Data = rpc_resp
	resp.ErrCode = pb.Error_EN_Success
	resp.Message = "题目提交成功"
	ctx.JSON(http.StatusOK, resp)
}

// 处理获取提交结果
func HandleGetSubmitResult(ctx *gin.Context) {}

// 处理创建题目信息
func HandleCreateProblem(ctx *gin.Context) {
	form, ret := validate(ctx, define.CreateProblemForm{})
	if !ret {
		return
	}
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	// 调用problem服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.CreateProblemRequest{
		Title:       form.Title,
		Description: form.Description,
		Level:       form.Level,
		Tags:        form.Tags,
	}
	rpc_resp, err := client.CreateProblem(context.Background(), req)
	if err != nil {
		logrus.Errorf("problem服务创建题目失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "创建题目失败"
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.ErrCode = pb.Error_EN_Success
	resp.Data = rpc_resp
	resp.Message = "创建题目成功"
	ctx.JSON(http.StatusOK, resp)
}

// 处理题目配置文件上传
func HandleUploadConfig(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 获取元数据
	problemIdStr := ctx.PostForm("problem_id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil || problemId <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "无效的 problem_id"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 获取文件
	fileHeader, err := ctx.FormFile("config_file")
	if err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "需要config_file字段"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	if filepath.Ext(fileHeader.Filename) != ".json" {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "文件格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 限制文件大小（默认 8MB，可调整）
	maxSize := int64(8 << 20) // 8MB
	if fileHeader.Size > maxSize {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "文件大小超过限制"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "文件打开失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 创建gRPC流
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	stream, err := client.UploadConfig(ctx)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
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
			resp.ErrCode = pb.Error_EN_ServiceBusy
			resp.Message = "服务器错误"
			ctx.JSON(http.StatusInternalServerError, resp)
			return
		}
		// 发送分片到Problem服务
		err = stream.Send(&pb.UploadConfigFileChunk{
			Content:   buffer[:n],
			ProblemId: problemId,
		})
		if err != nil {
			logrus.Errorf("file send error:%s", err.Error())
			_ = stream.CloseSend()
			resp.ErrCode = pb.Error_EN_ServiceBusy
			resp.Message = "服务器错误"
			ctx.JSON(500, resp)
			return
		}
	}
	// 获取Problem服务响应
	rpc_resp, err := stream.CloseAndRecv()
	if err != nil {
		logrus.Errorf("file save error:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "文件保存失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.ErrCode = pb.Error_EN_Success
	resp.Message = "文件上传成功"
	resp.Data = rpc_resp
	ctx.JSON(http.StatusOK, resp)
}

// 处理发布题目
func HandlePublishProblem(ctx *gin.Context) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 获取元数据
	problemIdStr := ctx.PostForm("problem_id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil || problemId <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "无效的 problem_id"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 调用problem服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.PublishProblemRequest{
		Id: problemId,
	}
	rpc_resp, err := client.PublishProblem(ctx, req)
	if err != nil {
		logrus.Errorf("PublishProblem Failed: %s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Data = rpc_resp
	resp.Message = "发布成功"
	ctx.JSON(http.StatusOK, resp)
}

// 处理删除题目
func HandleDeleteProblem(ctx *gin.Context) {}

// 处理更新题目信息
func HandleUpdateProblem(ctx *gin.Context) {}
