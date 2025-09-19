package biz

import (
	"oj-server/module/db"
)

// 仓库接口由data层去实现
type RecordRepo interface {
	// 查询提交记录列表
	QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []db.SubmitRecord, error)
	// 查询提交记录
	QuerySubmitRecord(id int64) (*db.SubmitRecord, error)
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
