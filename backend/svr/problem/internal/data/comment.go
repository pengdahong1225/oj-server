package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/db"
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
		logrus.Errorf("query problem failed: %s", result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment obj[%d] not found", id)
		return false
	}
	return true
}

func (cr *CommentRepo) SaveRootComment(comment *db.Comment) {
	err := cr.db_.Transaction(func(tx *gorm.DB) error {
		// 插入评论
		result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
		if result.Error != nil {
			return fmt.Errorf("insert comment failed: %s", result.Error.Error())
		}

		// 更新obj的总评论数
		if !updateObjCommentCount(tx, comment.ObjId, 1) {
			return fmt.Errorf("update obj comment count failed")
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("save root comment failed: %s", err.Error())
	}
}
func (cr *CommentRepo) SaveChildComment(comment *db.Comment) {
	err := cr.db_.Transaction(func(tx *gorm.DB) error {
		if !assertRoot(tx, comment.RootCommentId, comment.RootId) {
			return fmt.Errorf("assert root comment failed")
		}
		if comment.ReplyCommentId > 0 && comment.ReplyId > 0 && !assertReply(tx, comment.ReplyCommentId, comment.ReplyId) {
			return fmt.Errorf("assert reply comment failed")
		}

		// 插入评论
		result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
		if result.Error != nil {
			return fmt.Errorf("insert comment failed: %s", result.Error.Error())
		}

		// 更新obj的评论数
		if !updateObjCommentCount(tx, comment.ObjId, 1) {
			return fmt.Errorf("update obj comment count failed")
		}
		// 更新root的子评论数
		if !updateRootCommentChildCount(tx, comment.RootCommentId, 1) {
			return fmt.Errorf("update root comment child count failed")
		}
		// 更新reply的回复数
		if !updateReplyCommentReplyCount(tx, comment.ReplyCommentId, 1) {
			return fmt.Errorf("update reply comment reply count failed")
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("save child comment failed: %s", err.Error())
	}
}
func (cr *CommentRepo) QueryRootComment(objId int64, page, pageSize int) (int64, []db.Comment, error) {
	/*
		select * from comment
		where obj_id = ? and is_root = 1
		order by like_count desc
		offset off_set
		limit pageSize;
	*/
	var count int64
	result := cr.db_.Model(&db.Comment{}).Where("obj_id = ?", objId).Where("is_root = ?", 1).Count(&count)
	if result.Error != nil {
		logrus.Errorf("query comment count failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment count failed")
	}
	if count == 0 {
		return 0, nil, nil
	}

	offSet := (page - 1) * pageSize
	var comments []db.Comment
	result = cr.db_.Where("obj_id = ?", objId).Where("is_root = ?", 1).
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

func (cr *CommentRepo) QueryChildComment(objId, rootId, rootCommentId int64, cursor int32) (int64, []db.Comment, error) {
	var count int64
	result := cr.db_.Model(&db.Comment{}).
		Where("obj_id = ?", objId).
		Where("is_root = ?", 0).
		Where("root_id = ?", rootId).
		Where("root_comment_id = ?", rootCommentId).
		Count(&count)
	if result.Error != nil {
		logrus.Errorf("query comment count failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment count failed")
	}
	if count == 0 {
		return 0, nil, nil
	}

	var comments []db.Comment
	result = cr.db_.Where("obj_id = ?", objId).
		Where("is_root = ?", 0).
		Where("root_id = ?", rootId).
		Where("root_comment_id = ?", rootCommentId).
		Order("id asc").
		Where("id > ?", cursor).
		Limit(5).
		Find(&comments)
	if result.Error != nil {
		logrus.Errorf("query comment failed: %s", result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "query comment failed")
	}

	return count, comments, nil
}

// 根评论是否存在，校验root_id和root_comment_id
func assertRoot(tx *gorm.DB, rootCommentId int64, rootId int64) bool {
	var (
		c db.Comment
	)
	result := tx.Where("id = ?", rootCommentId).Where("user_id = ?", rootId).Find(&c)
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
func (cr *CommentRepo) CommentLike(objId, commentId int64) {
	result := cr.db_.Model(&db.Comment{}).
		Where("obj_id = ?", objId).
		Where("id = ?", commentId).
		Update("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		logrus.Errorf("update comment like count failed: %s", result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] like update failed", commentId)
		return
	}
}

// 回复评论是否存在，校验reply_id和reply_comment_id
func assertReply(tx *gorm.DB, replyCommentId int64, replyId int64) bool {
	var (
		c db.Comment
	)
	result := tx.Where("id = ?", replyCommentId).Where("user_id = ?", replyId).Find(&c)
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

// 更新obj的评论数
func updateObjCommentCount(tx *gorm.DB, id int64, diff int64) bool {
	/*
		select comment_count from problem
		where id = ?;
	*/
	obj := db.Problem{}
	result := tx.Select("comment_count").Where("id = ?", id).Find(&obj)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("obj[%d] not found", id)
		return false
	}
	/*
		update problem set comment_count = comment_count + diff
		where id = ?;
	*/
	count := obj.CommentCount + diff
	if count < 0 {
		count = 0
	}
	result = tx.Model(&obj).Where("id = ?", id).Update("comment_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("obj[%d] comment update failed", id)
		return false
	}
	return true
}
func updateRootCommentChildCount(tx *gorm.DB, id int64, diff int) bool {
	/*
		select child_count from comment
		where id = ?;
	*/
	comment := db.Comment{}
	result := tx.Select("child_count").Where("id = ?", id).Find(&comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] not found", id)
		return false
	}
	/*
		update comment set child_count = child_count + diff
		where id = ?;
	*/
	count := comment.ChildCount + diff
	if count < 0 {
		count = 0
	}
	result = tx.Model(&comment).Where("id = ?", id).Update("child_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] child_count update failed", id)
		return false
	}
	return true
}
func updateReplyCommentReplyCount(tx *gorm.DB, id int64, diff int) bool {
	/*
		select reply_count from comment
		where id = ?;
	*/
	comment := db.Comment{}
	result := tx.Select("reply_count").Where("id = ?", id).Find(&comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] not found", id)
		return false
	}
	/*
		update comment set reply_count = reply_count + diff
		where id = ?;
	*/
	count := comment.ReplyCount + diff
	if count < 0 {
		count = 0
	}
	result = tx.Model(&comment).Where("id=?", id).Update("reply_count", count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] reply_count update failed", id)
		return false
	}
	return true
}
