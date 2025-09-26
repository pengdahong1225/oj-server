package biz

import (
	"oj-server/module/db"
	"oj-server/proto/pb"
)

// 仓库接口由data层去实现
type Repo interface {
	// 查询用户信息
	QueryUserInfo(uid int64) (*db.UserInfo, error)

	// 查询题目信息
	QueryProblemData(id int64) (*db.Problem, error)

	// 1.更新提交记录表
	// 2.更新解题记录表
	// 3.更新统计表
	UpdateUserSubmitRecord(record *db.SubmitRecord, level int32) error

	// 查询用户 当日和当月 的解题数量
	QueryUserAcceptCount(uid int64) (int64, int64, error)

	// 解锁
	UnLock(key string) error
	// 设置任务状态
	SetTaskState(taskId string, state int) error
	// 设置判题结果
	SetTaskResult(taskId, result string) error

	// 更新排行榜
	UpdateLeaderboard(targetKey string, lb *pb.LeaderboardUserInfo) error
}

type JudgeUseCase struct {
	repo Repo
}

func NewJudgeUseCase(repo Repo) *JudgeUseCase {
	return &JudgeUseCase{
		repo: repo,
	}
}
func (jc *JudgeUseCase) QueryUserInfo(uid int64) (*db.UserInfo, error) {
	return jc.repo.QueryUserInfo(uid)
}

func (jc *JudgeUseCase) QueryProblemData(id int64) (*db.Problem, error) {
	return jc.repo.QueryProblemData(id)
}
func (jc *JudgeUseCase) UpdateUserSubmitRecord(record *db.SubmitRecord, level int32) error {
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

func (jc *JudgeUseCase) UpdateLeaderboard(targetKey string, lb *pb.LeaderboardUserInfo) error {
	return jc.repo.UpdateLeaderboard(targetKey, lb)
}

func (jc *JudgeUseCase) QueryUserAcceptCount(uid int64) (int64, int64, error) {
	return jc.repo.QueryUserAcceptCount(uid)
}
