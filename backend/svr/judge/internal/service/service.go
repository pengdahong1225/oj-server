package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/module/db"
	"oj-server/module/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"oj-server/svr/judge/internal/data"
	"oj-server/svr/judge/internal/define"
	"oj-server/svr/judge/internal/processor"
	"oj-server/utils"
	"os"
	"time"
)

var (
	uc *biz.JudgeUseCase
)

func Init() error {
	repo, err := data.NewRepo()
	if err != nil {
		return err
	}
	uc = biz.NewJudgeUseCase(repo)

	return nil
}

func Handle(form *pb.SubmitForm) {
	var results []*pb.PBResult

	// 预处理，构造上下文参数
	ok, param := preAction(form)
	if !ok {
		logrus.Errorln("预处理失败")
		results = append(results, &pb.PBResult{
			Status: define.InternalError,
			ErrMsg: "任务预处理失败",
		})
	} else {
		p := processor.NewProcessor(form.Lang, uc)
		if p == nil {
			return
		}

		start := time.Now()
		results = p.HandleJudgeTask(form, param)
		duration := time.Now().Sub(start).Milliseconds()
		logrus.Infof("---Handle--- uid:%d, problemID:%d, total-cost:%d ms", form.Uid, form.ProblemId, duration)
	}

	analyzeResult(param, results)
	// 记录结果
	results_data, err := json.Marshal(results)
	if err != nil {
		logrus.Errorf("结果序列化错误, err=%s", err.Error())
		return
	}
	saveResult(param, results_data)
}

func preAction(form *pb.SubmitForm) (bool, *define.Param) {
	param := &define.Param{}

	// 拉取题目信息
	problem, err := uc.QueryProblemData(form.ProblemId)
	if err != nil {
		logrus.Errorf("无法拉取题目[%d]信息, err=%s", form.ProblemId, err.Error())
		return false, nil
	}

	// 读取题目配置文件
	cfg_path := fmt.Sprintf("%s/%d.json", global.ProblemConfigPath, form.ProblemId)
	if _, err := os.Stat(cfg_path); os.IsNotExist(err) {
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
	param.ProblemData = utils.Transform(problem)
	param.Code = form.Code
	param.Language = form.Lang
	param.ProblemConfig = &cfg
	return true, param
}

// 分析结果
func analyzeResult(param *define.Param, results []*pb.PBResult) {
	param.Accepted = true
	param.Message = define.Accepted
	for _, result := range results {
		if result.Status != define.Accepted {
			param.Accepted = false
			param.Message = result.Content
			break
		}
	}
}

func saveResult(param *define.Param, data []byte) {
	// 保存本次提交结果
	taskId := fmt.Sprintf("%d_%d", param.Uid, param.ProblemData.Id)
	err := uc.SetTaskResult(taskId, param.Message)
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
	err = uc.UpdateUserSubmitRecord(record, param.ProblemData.Level)
	if err != nil {
		logrus.Errorf("更新数据库提交记录失败, err=%s", err.Error())
	}
}
