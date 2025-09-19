package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/db"
	"oj-server/module/proto/pb"
)

type CommentRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewCommentRepo() (*CommentRepo, error) {
	mysql_cfg := configManager.AppConf.MysqlCfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	db_, err := db.NewMysqlCli(dsn, global.LogPath)
	if err != nil {
		return nil, err
	}

	redis_cfg := configManager.AppConf.RedisCfg
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_, err := db.NewRedisCli(dsn)
	if err != nil {
		return nil, err
	}

	return &CommentRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (cr *CommentRepo) AssertObj(id int64) bool {
	var (
		p db.Problem
	)
	result := cr.db_.Where("id = ?", id).Find(&p)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment obj[%d] not found", id)
		return false
	}
	return true
}
func (cr *CommentRepo) AssertRoot(rootCommentId int64, rootId int64) bool {
	var (
		c db.Comment
	)
	result := cr.db_.Where("id = ?", rootCommentId).Where("user_id = ?", rootId).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("root comment[%d] not found", rootCommentId)
		return false
	}
	return true
}
func (cr *CommentRepo) AssertReply(replyCommentId int64, replyId int64) bool {
	var (
		c db.Comment
	)
	result := cr.db_.Where("id = ?", replyCommentId).Where("user_id = ?", replyId).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("reply comment[%d] not found", replyCommentId)
		return false
	}
	return true
}
func (cr *CommentRepo) SaveRootComment(pbComment *pb.Comment) {

}
func (cr *CommentRepo) SaveChildComment(pbComment *pb.Comment) {

}
func (cr *CommentRepo) QueryRootComment(page, pageSize int) (int32, []pb.Comment, error) {
}
