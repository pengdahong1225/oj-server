package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/module/db"
	"oj-server/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"oj-server/svr/judge/internal/processor"
	"oj-server/utils"
	"os"
	"time"
)

func (s *JudgeService) Handle(form *pb.SubmitForm) {
	var results []*pb.PBResult

	// 预处理，构造上下文参数
	ok, param := s.preAction(form)
	if !ok {
		logrus.Errorln("预处理失败")
		results = append(results, &pb.PBResult{
			Status: biz.InternalError,
			ErrMsg: "任务预处理失败",
		})
	} else {
		// 获取处理器
		p, err := processor.NewProcessor(form.Lang, s.uc)
		if err != nil {
			logrus.Errorf("无法获取处理器, err=%s", err.Error())
			return
		}

		start := time.Now()
		results = p.HandleJudgeTask(form, param)
		duration := time.Now().Sub(start).Milliseconds()
		logrus.Infof("---Handle Judge Task--- uid:%d, problemID:%d, total-cost:%d ms", form.Uid, form.ProblemId, duration)
	}

	s.analyzeResult(param, results)
	// 记录结果
	results_data, err := json.Marshal(results)
	if err != nil {
		logrus.Errorf("结果序列化错误, err=%s", err.Error())
		return
	}
	s.saveResult(param, results_data)
}

func (s *JudgeService) preAction(form *pb.SubmitForm) (bool, *biz.Param) {
	param := &biz.Param{}

	// 拉取题目信息
	problem, err := s.uc.QueryProblemData(form.ProblemId)
	if err != nil {
		logrus.Errorf("无法拉取题目[%d]信息, err=%s", form.ProblemId, err.Error())
		return false, nil
	}

	// 读取题目配置文件
	cfg_path := fmt.Sprintf("%s/%d.json", global.ProblemConfigPath, form.ProblemId)
	if _, err = os.Stat(cfg_path); os.IsNotExist(err) {
		logrus.Errorf("题目[%d]配置文件不存在, err=%s", form.ProblemId, err.Error())
		return false, nil
	}
	file_data, err := os.ReadFile(cfg_path)
	if err != nil {
		logrus.Errorf("无法读取题目[%d]配置文件, err=%s", form.ProblemId, err.Error())
		return false, nil
	}
	var cfg pb.ProblemConfig
	err = json.Unmarshal(file_data, &cfg)
	if err != nil {
		logrus.Errorf("解析题目[%d]配置文件错误, err=%s", form.ProblemId, err.Error())
		return false, nil
	}

	param.Uid = form.Uid
	param.UserName = form.UserName
	param.ProblemData = &pb.Problem{
		Id:          problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        utils.SplitStringWithX(string(problem.Tags), "#"),
		CreateBy:    problem.CreateBy,
		CreateAt:    problem.CreateAt.String(),
		Status:      problem.Status,
	}
	param.Code = form.Code
	param.Language = form.Lang
	param.ProblemConfig = &cfg
	return true, param
}

// 分析结果
func (s *JudgeService) analyzeResult(param *biz.Param, results []*pb.PBResult) {
	param.Accepted = true
	param.Message = biz.Accepted
	for _, result := range results {
		if result.Status != biz.Accepted {
			param.Accepted = false
			param.Message = result.Content
			break
		}
	}
}

func (s *JudgeService) saveResult(param *biz.Param, data []byte) {
	// 保存本次提交结果
	taskId := fmt.Sprintf("%d_%d", param.Uid, param.ProblemData.Id)
	err := s.uc.SetTaskResult(taskId, param.Message)
	if err != nil {
		logrus.Errorf("保存判题结果失败, err=%s", err.Error())
	}
	// 更新数据库
	record := &db.SubmitRecord{
		Uid:         param.Uid,
		UserName:    param.UserName,
		ProblemID:   param.ProblemData.Id,
		ProblemName: param.ProblemData.Title,
		Status:      param.Message,
		Code:        param.Code,
		Result:      data,
		Lang:        param.Language,
	}
	err = s.uc.UpdateUserSubmitRecord(record, param.ProblemData.Level)
	if err != nil {
		logrus.Errorf("更新数据库提交记录失败, err=%s", err.Error())
	}
}
