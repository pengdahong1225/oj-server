package comment

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/sirupsen/logrus"
)

type Querier struct {
}

func (c *Querier) onRootComment(objId int64, cursor int64) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 1 and id > cursor
		order by like_count desc
		limit 10;
	*/
	var comments []mysql.Comment
	db := mysql.DBSession
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 1).Where("obj_id > ?", cursor).Order("like_count desc").Limit(10).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}
func (c *Querier) onChildComment(objId int64, cursor int64) []mysql.Comment {
	/*
		select * from comment
		where obj_id = ? and is_root = 0 and id > cursor
		limit 5;
	*/
	var comments []mysql.Comment
	db := mysql.DBSession
	result := db.Where("obj_id = ?", objId).Where("is_root = ?", 0).Where("obj_id > ?", cursor).Limit(5).Find(&comments)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
	}
	return comments
}
