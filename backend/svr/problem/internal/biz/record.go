package biz

import (
	"oj-server/pkg/proto/pb"
	"oj-server/svr/problem/internal/model"
	"time"
)

// 仓库接口由data层去实现
type RecordRepo interface {
	// 数据库接口
	// 查询提交记录列表
	QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []model.SubmitRecord, error)
	// 查询提交记录
	QuerySubmitRecord(id int64) (*model.SubmitRecord, error)
	// 根据uid查询解题统计信息
	QueryStatistics(uid int64) (*model.Statistics, error)

	// 查询统计表 -- 月榜
	QueryMonthAccomplishLeaderboard(limit int, period string) ([]*pb.LeaderboardUserInfo, error)
	// 查询解题表 -- 日榜
	QueryDailyAccomplishLeaderboard(limit int) ([]*pb.LeaderboardUserInfo, error)

	// redis接口
	// 获取排行榜上次更新时间
	QueryLeaderboardLastUpdate() (int64, error)
	// 更新排行榜上次更新时间
	UpdateLeaderboardLastUpdate(time int64) error
	// 同步排行榜
	SynchronizeLeaderboard(lb_list []*pb.LeaderboardUserInfo, leaderboardKey string, leaderboardKeyTTL time.Duration) error
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
func (rc *RecordUseCase) QueryLeaderboardLastUpdate() (int64, error) {
	return rc.repo.QueryLeaderboardLastUpdate()
}
func (rc *RecordUseCase) UpdateLeaderboardLastUpdate(time int64) error {
	return rc.repo.UpdateLeaderboardLastUpdate(time)
}
func (rc *RecordUseCase) QueryStatistics(uid int64) (*model.Statistics, error) {
	return rc.repo.QueryStatistics(uid)
}
func (rc *RecordUseCase) QueryMonthAccomplishLeaderboard(limit int, period string) ([]*pb.LeaderboardUserInfo, error) {
	return rc.repo.QueryMonthAccomplishLeaderboard(limit, period)
}
func (rc *RecordUseCase) QueryDailyAccomplishLeaderboard(limit int) ([]*pb.LeaderboardUserInfo, error) {
	return rc.repo.QueryDailyAccomplishLeaderboard(limit)
}
func (rc *RecordUseCase) SynchronizeLeaderboard(lb_list []*pb.LeaderboardUserInfo, leaderboardKey string, leaderboardKeyTTL time.Duration) error {
	return rc.repo.SynchronizeLeaderboard(lb_list, leaderboardKey, leaderboardKeyTTL)
}
