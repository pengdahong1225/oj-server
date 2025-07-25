package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/judge/internal/define"
	"oj-server/app/judge/internal/processor"
	"oj-server/app/judge/internal/respository/cache"
	"oj-server/app/judge/internal/respository/domain"
	"oj-server/global"
	"oj-server/module/model"
	"oj-server/proto/pb"
	"os"
	"time"
)

var (
	db_ *domain.MysqlDB
)

func Init() error {
	var err error
	db_, err = domain.NewMysqlDB()
	return err
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
		p := processor.NewProcessor(form.Lang)

		start := time.Now()
		results = p.HandleJudgeTask(form, param)
		duration := time.Now().Sub(start).Milliseconds()
		logrus.Infof("---Handle--- uid:%d, problemID:%d, total-cost:%d ms", form.Uid, form.ProblemId, duration)
	}

	analyzeResult(param, results)
	// 记录结果
	data, err := json.Marshal(results)
	if err != nil {
		logrus.Errorf("结果序列化错误, err=%s", err.Error())
		return
	}
	saveResult(param, data)
}

func preAction(form *pb.SubmitForm) (bool, *define.Param) {
	param := &define.Param{}

	// 拉取题目信息
	problem, err := db_.QueryProblemData(form.ProblemId)
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
	param.ProblemData = problem.Transform()
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
	err := cache.SetJudgeResult(taskId, param.Message, global.JudgeResultExpired)
	if err != nil {
		logrus.Errorf("保存判题结果失败, err=%s", err.Error())
	}
	// 更新数据库
	record := &model.SubmitRecord{
		Uid:         param.Uid,
		UserName:    param.UserName,
		ProblemID:   param.ProblemData.Id,
		ProblemName: param.ProblemData.Title,
		Status:      param.Message,
		Code:        param.Code,
		Result:      data,
		Lang:        param.Language,
	}
	err = db_.UpdateUserSubmitRecord(record, param.ProblemData.Level)
	if err != nil {
		logrus.Errorf("更新数据库提交记录失败, err=%s", err.Error())
	}
}
