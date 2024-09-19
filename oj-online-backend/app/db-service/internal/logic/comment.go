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
	// 校验
	checker := commentChecker{}
	if !checker.assertObj(request.ObjId) {
		logrus.Errorf("obj[%d] assert failed", request.ObjId)
		return nil, QueryFailed
	}
	if request.RootId > 0 && request.RootCommentId > 0 && !checker.assertRoot(request.RootCommentId, request.RootId) {
		logrus.Errorf("root comment[%d] assert failed", request.RootCommentId)
		return nil, QueryFailed
	}
	if request.ReplyId > 0 && request.ReplyCommentId > 0 && !checker.assertReply(request.ReplyCommentId, request.ReplyId) {
		logrus.Errorf("reply comment[%d] assert failed", request.ReplyCommentId)
		return nil, QueryFailed
	}

	response := &pb.QueryCommentResponse{}
	cursor := request.Cursor
	if cursor < 0 {
		cursor = 0
	}
	if request.RootId > 0 && request.RootCommentId > 0 {
		comments := receiver.onRootComment(request.ObjId, request.Cursor)
		for _, comment := range comments {
			response.Data = append(response.Data, translateComment(&comment))
		}
		return response, nil
	} else {
		comments := receiver.onChildComment(request.ObjId, request.Cursor)
		for _, comment := range comments {
			response.Data = append(response.Data, translateComment(&comment))
		}
		return response, nil
	}
}
func (receiver *DBServiceServer) onRootComment(objId int64, cursor int64) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 1 and id > cursor
		order by like_count desc
		limit 10;
	*/
	var comments []mysql.Comment
	db := mysql.Instance()
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 1).Where("obj_id > ?", cursor).Order("like_count desc").Limit(10).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}
func (receiver *DBServiceServer) onChildComment(objId int64, cursor int64) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 0 and id > cursor
		limit 5;
	*/
	var comments []mysql.Comment
	db := mysql.Instance()
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 0).Where("obj_id > ?", cursor).Limit(5).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}

func (receiver *DBServiceServer) DeleteComment(ctx context.Context, request *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}

type CommentSaveHandler struct{}

func (receiver CommentSaveHandler) SaveComment(pbComment *pb.Comment) {
	checker := commentChecker{}
	if !checker.assertObj(pbComment.ObjId) {
		logrus.Errorf("obj[%d] assert failed", pbComment.ObjId)
		return
	}
	if pbComment.IsRoot > 0 {
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
	checker := commentChecker{}
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

// assert
type commentChecker struct {
}

// obj是否存在
func (receiver commentChecker) assertObj(id int64) bool {
	db := mysql.Instance()

	var p mysql.Problem
	result := db.Where("id = ?", id).Find(&p)
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

// 根评论是否存在，校验root_id和root_comment_id
func (receiver commentChecker) assertRoot(rootCommentId int64, rootId int64) bool {
	db := mysql.Instance()
	var c mysql.Comment
	result := db.Where("id = ?", rootCommentId).Where("user_id = ?", rootId).Find(&c)
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

// 回复评论是否存在，校验reply_id和reply_comment_id
func (receiver commentChecker) assertReply(replyCommentId int64, replyId int64) bool {
	db := mysql.Instance()
	var c mysql.Comment
	result := db.Where("id = ?", replyCommentId).Where("user_id = ?", replyId).Find(&c)
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

func translateComment(comment *mysql.Comment) *pb.Comment {
	return &pb.Comment{
		Id:            int64(comment.ID),
		ObjId:         comment.ObjId,
		UserId:        comment.UserId,
		UserName:      comment.UserName,
		UserAvatarUrl: comment.UserAvatarUrl,
		Content:       comment.Content,
		Status:        int32(comment.Status),
		LikeCount:     int32(comment.LikeCount),
		ReplyCount:    int32(comment.ReplyCount),
		ChildCount:    int32(comment.ChildCount),
		Stamp:         comment.Stamp,

		IsRoot:         int32(comment.IsRoot),
		RootId:         comment.RootId,
		RootCommentId:  comment.RootCommentId,
		ReplyId:        comment.ReplyId,
		ReplyCommentId: comment.ReplyCommentId,
	}
}
