package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"oj-server/svr/judge/internal/processor"
	"os"
	"time"
)

func (s *JudgeService) Handle(task *pb.JudgeSubmission) {
	judgeResult := &pb.JudgeResult{
		Uid:       task.Uid,
		ProblemId: task.ProblemId,
	}

	// 预处理，构造上下文参数
	ok, param := s.preAction(task)
	if !ok {
		logrus.Errorf("预处理失败")
		judgeResult.Items = append(judgeResult.Items, &pb.JudgeResultItem{
			Status: biz.InternalError,
			ErrMsg: "任务预处理失败",
		})
	} else {
		// 获取处理器
		p, err := processor.NewProcessor(task.Lang)
		if err != nil {
			logrus.Errorf("无法获取处理器, err=%s", err.Error())
			return
		}

		start := time.Now()
		results := p.HandleJudgeTask(task, param)
		duration := time.Now().Sub(start).Milliseconds()
		logrus.Infof("---Handle Judge Task--- uid:%d, problemID:%d, total-cost:%d ms",
			task.Uid, task.ProblemId, duration)

		judgeResult.Items = append(judgeResult.Items, results...)
	}

	// 处理结果
	param.Accepted = true
	param.Message = biz.Accepted
	for _, item := range judgeResult.Items {
		if item.Status != biz.Accepted {
			param.Accepted = false
			param.Message = item.Content
			break
		}
	}
	judgeResult.Accepted = param.Accepted
	judgeResult.Message = param.Message
	judgeResult.Code = param.Code
	judgeResult.Lang = param.Language
	judgeResult.TaskId = param.TaskId

	// 推送到队列
	err := s.sendJudgeResult2MQ(judgeResult)
	if err != nil {
		logrus.Errorln(err)
	}
}

func (s *JudgeService) preAction(task *pb.JudgeSubmission) (bool, *biz.Param) {
	param := &biz.Param{}

	param.Code = task.Code
	param.Language = task.Lang
	param.TaskId = task.TaskId

	// 读取题目配置文件
	cfg_path := task.ConfigUrl
	if _, err := os.Stat(cfg_path); os.IsNotExist(err) {
		logrus.Errorf("题目[%d]配置文件不存在, err=%s", task.ProblemId, err.Error())
		return false, nil
	}
	file_data, err := os.ReadFile(cfg_path)
	if err != nil {
		logrus.Errorf("无法读取题目[%d]配置文件, err=%s", task.ProblemId, err.Error())
		return false, nil
	}
	cfg := new(pb.ProblemConfig)
	if err = json.Unmarshal(file_data, cfg); err != nil {
		logrus.Errorf("解析题目[%d]配置文件错误, err=%s", task.ProblemId, err.Error())
		return false, nil
	}

	param.ProblemConfig = cfg

	return true, param
}

func (s *JudgeService) sendJudgeResult2MQ(result *pb.JudgeResult) error {
	data, err := proto.Marshal(result)
	if err != nil {
		return err
	}
	return s.result_producer.Publish(data)
}
