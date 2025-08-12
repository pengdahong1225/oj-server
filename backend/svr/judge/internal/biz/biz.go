package biz

import (
	model2 "oj-server/module/db/model"
)

// 仓库接口由data层去实现
type Repo interface {
	// 查询题目信息
	QueryProblemData(id int64) (*model2.Problem, error)

	// 1.更新用户提交记录表
	// 2.更新(用户,题目)AC表
	// 3.更新用户解题情况统计表
	UpdateUserSubmitRecord(record *model2.SubmitRecord, level int32) error

	// 解锁
	UnLock(key string) error
	// 设置任务状态
	SetTaskState(taskId string, state int) error
	// 设置判题结果
	SetTaskResult(taskId, result string) error
}

type JudgeUseCase struct {
	repo Repo
}

func NewJudgeUseCase(repo Repo) *JudgeUseCase {
	return &JudgeUseCase{
		repo: repo,
	}
}
func (jc *JudgeUseCase) QueryProblemData(id int64) (*model2.Problem, error) {
	return jc.repo.QueryProblemData(id)
}
func (jc *JudgeUseCase) UpdateUserSubmitRecord(record *model2.SubmitRecord, level int32) error {
	return jc.repo.UpdateUserSubmitRecord(record, level)
}
func (jc *JudgeUseCase) SetTaskState(taskId string, state int) error {
	return jc.repo.SetTaskState(taskId, state)
}
func (jc *JudgeUseCase) UnLock(key string) error {
	return jc.repo.UnLock(key)
}
func (jc *JudgeUseCase) SetTaskResult(taskId, result string) error {
	return jc.repo.SetTaskResult(taskId, result)
}
