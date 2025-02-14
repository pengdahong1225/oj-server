package comment

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Saver struct {
}

func (s *Saver) onRootComment(pbComment *pb.Comment) {
	comment := mysql.NewComment()
	comment.ObjId = pbComment.ObjId
	comment.UserId = pbComment.UserId
	comment.UserName = pbComment.UserName
	comment.UserAvatarUrl = pbComment.UserAvatarUrl
	comment.Content = pbComment.Content
	comment.PubRegion = pbComment.PubRegion
	comment.IsRoot = 1
	comment.RootId = 0
	comment.RootCommentId = 0
	comment.ReplyId = 0
	comment.ReplyCommentId = 0

	// 开启事务
	tx := mysql.DBSession.Begin()
	if tx.Error != nil {
		logrus.Errorln(tx.Error.Error())
		return
	}

	// 插入评论
	result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return
	}

	// 更新obj的总评论数
	if !s.updateObjCommentCount(tx, comment.ObjId, 1) {
		tx.Rollback()
		return
	}

	// 提交事务
	tx.Commit()
}
func (s *Saver) onChildComment(pbComment *pb.Comment) {
	checker := Checker{}
	if !checker.assertRoot(pbComment.RootCommentId, pbComment.RootId) {
		logrus.Errorf("root comment[%d] assert failed", pbComment.RootCommentId)
		return
	}
	if pbComment.ReplyId > 0 && pbComment.ReplyCommentId > 0 && !checker.assertReply(pbComment.ReplyCommentId, pbComment.ReplyId) {
		logrus.Errorf("reply comment[%d] assert failed", pbComment.ReplyCommentId)
		return
	}

	comment := mysql.NewComment()
	comment.ObjId = pbComment.ObjId
	comment.UserId = pbComment.UserId
	comment.UserName = pbComment.UserName
	comment.UserAvatarUrl = pbComment.UserAvatarUrl
	comment.Content = pbComment.Content
	comment.IsRoot = 0
	comment.RootId = pbComment.RootId
	comment.RootCommentId = pbComment.RootCommentId
	comment.ReplyId = pbComment.ReplyId
	comment.ReplyCommentId = pbComment.ReplyCommentId
	comment.ReplyUserName = pbComment.ReplyUserName

	tx := mysql.DBSession.Begin()
	if tx.Error != nil {
		logrus.Errorln(tx.Error.Error())
		tx.Rollback()
		return
	}
	// 插入评论
	if pbComment.ReplyId > 0 && pbComment.ReplyCommentId > 0 {
		comment.ReplyId = pbComment.ReplyId
		comment.ReplyCommentId = pbComment.ReplyCommentId
	} else {
		comment.ReplyId = 0
		comment.ReplyCommentId = 0
	}
	result := tx.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt").Create(comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return
	} else if result.RowsAffected == 0 {
		logrus.Errorf("comment[%d] insert failed", comment.ID)
		tx.Rollback()
		return
	}
	// 更新obj的总评论数
	if !s.updateObjCommentCount(tx, comment.ObjId, 1) {
		tx.Rollback()
		return
	}
	// 更新root的子评论数
	if !s.updateRootCommentChildCount(tx, comment.RootCommentId, 1) {
		tx.Rollback()
		return
	}
	// 更新reply的总回复数
	if !s.updateReplyCommentReplyCount(tx, comment.ReplyCommentId, 1) {
		tx.Rollback()
		return
	}
	// 提交事务
	tx.Commit()
}

func (s *Saver) updateObjCommentCount(tx *gorm.DB, id int64, diff int64) bool {
	/*
		select comment_count from problem
		where id = ?;
	*/
	obj := mysql.Problem{}
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
func (s *Saver) updateRootCommentChildCount(tx *gorm.DB, id int64, diff int) bool {
	/*
		select child_count from comment
		where id = ?;
	*/
	comment := mysql.Comment{}
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
func (s *Saver) updateReplyCommentReplyCount(tx *gorm.DB, id int64, diff int) bool {
	/*
		select reply_count from comment
		where id = ?;
	*/
	comment := mysql.Comment{}
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
