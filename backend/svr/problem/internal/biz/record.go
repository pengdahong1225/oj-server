package biz

import (
	"oj-server/module/db"
)

// 仓库接口由data层去实现
type RecordRepo interface {
	// 数据库接口
	// 查询提交记录列表
	QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []db.SubmitRecord, error)
	// 查询提交记录
	QuerySubmitRecord(id int64) (*db.SubmitRecord, error)
	// 根据uid查询解题统计信息
	QueryStatistics(uid int64) (*db.Statistics, error)

	// redis接口
	// 获取排行榜上次更新时间
	QueryLeaderboardLastUpdate() (int64, error)
	// 更新排行榜上次更新时间
	UpdateLeaderboardLastUpdate(time int64) error
}

type RecordUseCase struct {
	repo RecordRepo
}

func NewRecordUseCase(repo RecordRepo) *RecordUseCase {
	return &RecordUseCase{
		repo: repo,
	}
}

func (rc *RecordUseCase) QuerySubmitRecordList(uid int64, page, pageSize int) (int64, []db.SubmitRecord, error) {
	return rc.repo.QuerySubmitRecordList(uid, pageSize, page)
}
func (rc *RecordUseCase) QuerySubmitRecord(id int64) (*db.SubmitRecord, error) {
	return rc.repo.QuerySubmitRecord(id)
}
func (rc *RecordUseCase) QueryLeaderboardLastUpdate() (int64, error) {
	return rc.repo.QueryLeaderboardLastUpdate()
}
func (rc *RecordUseCase) UpdateLeaderboardLastUpdate(time int64) error {
	return rc.repo.UpdateLeaderboardLastUpdate(time)
}
func (rc *RecordUseCase) QueryStatistics(uid int64) (*db.Statistics, error) {
	return rc.repo.QueryStatistics(uid)
}
