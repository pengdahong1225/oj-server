package comment

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/sirupsen/logrus"
)

type Querier struct {
}

func (c *Querier) onRootComment(objId int64, page int32, pageSize int32) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 1
		order by like_count desc
		offset off_set
		limit pageSize;
	*/
	offSet := int((page - 1) * pageSize)

	var comments []mysql.Comment
	db := mysql.DBSession
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 1).Order("like_count desc").Offset(offSet).Limit(int(pageSize)).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}
func (c *Querier) onChildComment(objId, rootId, rootCommentId int64, cursor int32) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 0 and rootId=? and rootCommentId=? and id > cursor
		limit 5;
	*/
	var comments []mysql.Comment
	db := mysql.DBSession
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 0).Where("root_id = ?", rootId).Where("root_comment_id = ?", rootCommentId).Where("id > ?", cursor).Limit(5).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}
