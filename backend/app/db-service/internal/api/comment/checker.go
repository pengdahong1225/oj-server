package comment

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/sirupsen/logrus"
)

// Checker assert
type Checker struct {
}

// obj是否存在
func (receiver Checker) assertObj(id int64) bool {
	db := mysql.DBSession

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
func (receiver Checker) assertRoot(rootCommentId int64, rootId int64) bool {
	db := mysql.DBSession
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
func (receiver Checker) assertReply(replyCommentId int64, replyId int64) bool {
	db := mysql.DBSession
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
