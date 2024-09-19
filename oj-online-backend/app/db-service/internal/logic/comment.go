package logic

import (
	"context"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (receiver *DBServiceServer) QueryComment(ctx context.Context, request *pb.QueryCommentRequest) (*pb.QueryCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryComment not implemented")
}
func (receiver *DBServiceServer) DeleteComment(ctx context.Context, request *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}

type CommentSaveHandler struct{}

func (receiver CommentSaveHandler) SaveComment(pbComment *pb.Comment) {
	if !receiver.assertObj(pbComment) {
		logrus.Errorf("obj[%d] assert failed", pbComment.ObjId)
		return
	}
	if pbComment.IsRoot {
		receiver.onRootComment(pbComment) // 第一层
	} else {
		receiver.onChildComment(pbComment) // 第二层
	}
}
func (receiver CommentSaveHandler) onRootComment(pbComment *pb.Comment) {
	comment := mysql.NewComment()
	comment.ObjId = pbComment.ObjId
	comment.UserId = pbComment.UserId
	comment.UserName = pbComment.UserName
	comment.UserAvatarUrl = pbComment.UserAvatarUrl
	comment.Content = pbComment.Content
	comment.IsRoot = 1
	comment.RootId = 0
	comment.RootCommentId = 0
	comment.ReplyId = 0
	comment.ReplyCommentId = 0

	tx := mysql.Instance().Begin()
	if tx.Error != nil {
		logrus.Errorln(tx.Error.Error())
		return
	}

	// 插入评论
	result := tx.Create(comment)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return
	}

	// 更新obj的总评论数
	if !receiver.updateObjCommentCount(tx, comment.ObjId, 1) {
		tx.Rollback()
		return
	}

	// 提交事务
	tx.Commit()
}
func (receiver CommentSaveHandler) onChildComment(pbComment *pb.Comment) {
	if !receiver.assertRoot(pbComment) {
		logrus.Errorf("root comment[%d] assert failed", pbComment.RootCommentId)
		return
	}
	if pbComment.ReplyId > 0 && pbComment.ReplyCommentId > 0 && !receiver.assertReply(pbComment) {
		logrus.Errorf("reply comment[%d] assert failed", pbComment.ReplyCommentId)
		return
	}

	comment := mysql.NewComment()
	comment.ObjId = pbComment.ObjId
	comment.UserId = pbComment.UserId
	comment.UserName = pbComment.UserName
	comment.UserAvatarUrl = pbComment.UserAvatarUrl
	comment.Content = pbComment.Content
	comment.IsRoot = 1
	comment.RootId = pbComment.RootId
	comment.RootCommentId = pbComment.RootCommentId

	tx := mysql.Instance().Begin()
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
	result := tx.Create(comment)
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
	if !receiver.updateObjCommentCount(tx, comment.ObjId, 1) {
		tx.Rollback()
		return
	}
	// 更新root的子评论数
	if !receiver.updateRootCommentChildCount(tx, comment.RootCommentId, 1) {
		tx.Rollback()
		return
	}
	// 更新reply的总回复数
	if !receiver.updateReplyCommentReplyCount(tx, comment.ReplyCommentId, 1) {
		tx.Rollback()
		return
	}
	// 提交事务
	tx.Commit()
}

func (receiver CommentSaveHandler) updateObjCommentCount(tx *gorm.DB, id int64, diff int64) bool {
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
	result = tx.Model(&obj).Update("comment_count", count)
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
func (receiver CommentSaveHandler) updateRootCommentChildCount(tx *gorm.DB, id int64, diff int) bool {
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
	result = tx.Model(&comment).Update("child_count", count)
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
func (receiver CommentSaveHandler) updateReplyCommentReplyCount(tx *gorm.DB, id int64, diff int) bool {
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
	result = tx.Model(&comment).Update("reply_count", count)
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

// obj是否存在
func (receiver CommentSaveHandler) assertObj(pbComment *pb.Comment) bool {
	db := mysql.Instance()

	var p mysql.Problem
	result := db.Where("id = ?", pbComment.ObjId).Limit(1).Find(&p)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("comment obj[%d] not found", pbComment.ObjId)
		return false
	}
	return true
}

// 根评论是否存在，校验root_id和root_comment_id
func (receiver CommentSaveHandler) assertRoot(pbComment *pb.Comment) bool {
	db := mysql.Instance()
	var c mysql.Comment
	result := db.Where("id = ?", pbComment.RootCommentId).Where("user_id = ?", pbComment.UserId).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("root comment[%d] not found", pbComment.RootCommentId)
		return false
	}
	return true
}

// 回复评论是否存在，校验reply_id和reply_comment_id
func (receiver CommentSaveHandler) assertReply(pbComment *pb.Comment) bool {
	db := mysql.Instance()
	var c mysql.Comment
	result := db.Where("id = ?", pbComment.ReplyCommentId).Where("user_id = ?", pbComment.UserId).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		logrus.Errorf("reply comment[%d] not found", pbComment.ReplyCommentId)
		return false
	}
	return true
}
