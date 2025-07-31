package biz

import (
	"oj-server/module/db/model"
	"time"
)

// 仓库接口由data层去实现
type RecordRepo interface {
	// 查询提交记录列表
	QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []model.SubmitRecord, error)
	// 查询提交记录
	QuerySubmitRecord(id int64) (*model.SubmitRecord, error)

	// 加锁
	Lock(key string, ttl time.Duration) (bool, error)
	// 解锁
	UnLock(key string) error
}

type RecordUseCase struct {
	repo RecordRepo
}

func NewRecordUseCase(repo RecordRepo) *RecordUseCase {
	return &RecordUseCase{
		repo: repo,
	}
}

func (rc *RecordUseCase) QuerySubmitRecordList(uid int64, page, pageSize int) (int64, []model.SubmitRecord, error) {
	return rc.repo.QuerySubmitRecordList(uid, pageSize, page)
}
func (rc *RecordUseCase) QuerySubmitRecord(id int64) (*model.SubmitRecord, error) {
	return rc.repo.QuerySubmitRecord(id)
}
func (rc *RecordUseCase) Lock(key string, ttl time.Duration) (bool, error) {
	return rc.repo.Lock(key, ttl)
}
func (rc *RecordUseCase) UnLock(key string) error {
	return rc.repo.UnLock(key)
}
