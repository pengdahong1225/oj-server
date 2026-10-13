package data

import (
	"context"
	"fmt"
	"io"
	"log"
	"oj-server/common/consts"
	"oj-server/svr/problem/internal/config"
	"oj-server/svr/problem/internal/model"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CommentRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewCommentRepo() (*CommentRepo, error) {
	// mysql
	mysql_cfg := config.MyConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	timer := time.Now().Format("2006_01_02")
	filePath := fmt.Sprintf("%s/orm.%s.log", consts.LogPath, timer)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	writer := io.MultiWriter(os.Stdout, file)
	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})

	// redis
	redis_cfg := config.MyConfig.Redis
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_ := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb_.Ping(context.Background())
	if st.Err() != nil {
		return nil, st.Err()
	}

	return &CommentRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (cr *CommentRepo) AssertObj(ID int64) bool {
	var (
		p model.Problem
	)
	result := cr.db_.Where("ID = ?", ID).Find(&p)
	if result.Error != nil {
		logrus.Errorf("query problem failed: %s", result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment obj[%d] not found", ID)
		return false
	}
	return true
}

func (cr *CommentRepo) SaveRootComment(comment *model.Comment) {
	err := cr.db_.Transaction(func(tx *gorm.DB) error {
		// 插入评论
		result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
		if result.Error != nil {
			return fmt.Errorf("insert comment failed: %s", result.Error.Error())
		}

		// 更新obj的总评论数
		if !updateObjCommentCount(tx, comment.ObjID, 1) {
			return fmt.Errorf("update obj comment count failed")
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("save root comment failed: %s", err.Error())
	}
}
func (cr *CommentRepo) SaveChildComment(comment *model.Comment) {
	err := cr.db_.Transaction(func(tx *gorm.DB) error {
		if !assertRoot(tx, comment.RootCommentID, comment.RootID) {
			return fmt.Errorf("assert root comment failed")
		}
		if comment.ReplyCommentID > 0 && comment.ReplyID > 0 && !assertReply(tx, comment.ReplyCommentID, comment.ReplyID) {
			return fmt.Errorf("assert reply comment failed")
		}

		// 插入评论
		result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
		if result.Error != nil {
			return fmt.Errorf("insert comment failed: %s", result.Error.Error())
		}

		// 更新obj的评论数
		if !updateObjCommentCount(tx, comment.ObjID, 1) {
			return fmt.Errorf("update obj comment count failed")
		}
		// 更新root的子评论数
		if !updateRootCommentChildCount(tx, comment.RootCommentID, 1) {
			return fmt.Errorf("update root comment child count failed")
		}
		// 更新reply的回复数
		if !updateReplyCommentReplyCount(tx, comment.ReplyCommentID, 1) {
			return fmt.Errorf("update reply comment reply count failed")
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("save child comment failed: %s", err.Error())
	}
}
func (cr *CommentRepo) QueryRootComment(objID int64, page, pageSize int) (int64, []model.Comment, error) {
	/*
		select * from comment
		where obj_ID = ? and is_root = 1
		order by like_count desc
		offset off_set
		limit pageSize;
	*/
	var count int64
	result := cr.db_.Model(&model.Comment{}).Where("obj_ID = ?", objID).Where("is_root = ?", 1).Count(&count)
	if result.Error != nil {
		logrus.Errorf("query comment count failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment count failed")
	}
	if count == 0 {
		return 0, nil, nil
	}

	offSet := (page - 1) * pageSize
	var comments []model.Comment
	result = cr.db_.Where("obj_ID = ?", objID).Where("is_root = ?", 1).
		Order("like_count desc").
		Offset(offSet).
		Limit(pageSize).
		Find(&comments)
	if result.Error != nil {
		logrus.Errorf("query comment failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment failed")
	}

	return count, comments, nil
}

func (cr *CommentRepo) QueryChildComment(objID, rootID, rootCommentID int64, cursor int32) (int64, []model.Comment, error) {
	var count int64
	result := cr.db_.Model(&model.Comment{}).
		Where("obj_ID = ?", objID).
		Where("is_root = ?", 0).
		Where("root_ID = ?", rootID).
		Where("root_comment_ID = ?", rootCommentID).
		Count(&count)
	if result.Error != nil {
		logrus.Errorf("query comment count failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment count failed")
	}
	if count == 0 {
		return 0, nil, nil
	}

	var comments []model.Comment
	result = cr.db_.Where("obj_ID = ?", objID).
		Where("is_root = ?", 0).
		Where("root_ID = ?", rootID).
		Where("root_comment_ID = ?", rootCommentID).
		Order("ID asc").
		Where("ID > ?", cursor).
		Limit(5).
		Find(&comments)
	if result.Error != nil {
		logrus.Errorf("query comment failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment failed")
	}

	return count, comments, nil
}

// 根评论是否存在，校验root_ID和root_comment_ID
func assertRoot(tx *gorm.DB, rootCommentID int64, rootID int64) bool {
	var (
		c model.Comment
	)
	result := tx.Where("ID = ?", rootCommentID).Where("user_ID = ?", rootID).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("root comment[%d] not found", rootCommentID)
		return false
	}
	return true
}
func (cr *CommentRepo) CommentLike(objID, commentID int64) {
	result := cr.db_.Model(&model.Comment{}).
		Where("obj_ID = ?", objID).
		Where("ID = ?", commentID).
		Update("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		logrus.Errorf("update comment like count failed: %s", result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] like update failed", commentID)
		return
	}
}

// 回复评论是否存在，校验reply_ID和reply_comment_ID
func assertReply(tx *gorm.DB, replyCommentID int64, replyID int64) bool {
	var (
		c model.Comment
	)
	result := tx.Where("ID = ?", replyCommentID).Where("user_ID = ?", replyID).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("reply comment[%d] not found", replyCommentID)
		return false
	}
	return true
}

// 更新obj的评论数
func updateObjCommentCount(tx *gorm.DB, ID int64, diff int64) bool {
	/*
		select comment_count from problem
		where ID = ?;
	*/
	obj := model.Problem{}
	result := tx.Select("comment_count").Where("ID = ?", ID).Find(&obj)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("obj[%d] not found", ID)
		return false
	}
	/*
		update problem set comment_count = comment_count + diff
		where ID = ?;
	*/
	count := obj.CommentCount + diff
	if count < 0 {
		count = 0
	}
	result = tx.Model(&obj).Where("ID = ?", ID).Update("comment_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("obj[%d] comment update failed", ID)
		return false
	}
	return true
}
func updateRootCommentChildCount(tx *gorm.DB, ID int64, diff int) bool {
	/*
		select child_count from comment
		where ID = ?;
	*/
	comment := model.Comment{}
	result := tx.Select("child_count").Where("ID = ?", ID).Find(&comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] not found", ID)
		return false
	}
	/*
		update comment set child_count = child_count + diff
		where ID = ?;
	*/
	count := comment.ChildCount + int32(diff)
	if count < 0 {
		count = 0
	}
	result = tx.Model(&comment).Where("ID = ?", ID).Update("child_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] child_count update failed", ID)
		return false
	}
	return true
}
func updateReplyCommentReplyCount(tx *gorm.DB, ID int64, diff int) bool {
	/*
		select reply_count from comment
		where ID = ?;
	*/
	comment := model.Comment{}
	result := tx.Select("reply_count").Where("ID = ?", ID).Find(&comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] not found", ID)
		return false
	}
	/*
		update comment set reply_count = reply_count + diff
		where ID = ?;
	*/
	count := comment.ReplyCount + int32(diff)
	if count < 0 {
		count = 0
	}
	result = tx.Model(&comment).Where("ID=?", ID).Update("reply_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] reply_count update failed", ID)
		return false
	}
	return true
}
