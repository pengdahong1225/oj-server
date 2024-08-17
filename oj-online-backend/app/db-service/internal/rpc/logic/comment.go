package logic

import (
	mysql2 "github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
)

// 1.更新本条评论
// 2.更新回复评论的回复数
// 3.更新根评论（楼主）的child_count
func SaveComment(pbComment *pb.Comment) {
	db := mysql2.Instance()
	// obj_id目标评论区是否存在
	/*
		select count(*)
		from problem
		where id = ?;
	*/
	var c mysql2.Comment
	result := db.Where("id = ?", pbComment.ObjId).Limit(1).Find(&c)
	if result.Error != nil {
		logrus.Errorln(result.Error)
		return
	}
	if result.RowsAffected == 0 {
		logrus.Errorln("comment obj not found:", pbComment.ObjId)
		return
	}
	if !pbComment.IsRoot {
		// root_comment_id（根评论）是否存在
		// reply_comment_id（回复评论）是否存在
		/*
			select count(*)
			from comment
			where ObjId = ? and root_comment_id = ? and reply_comment_id = ?;
		*/

	}

	comment := mysql2.Comment{
		ObjId:         pbComment.ObjId,
		UserId:        pbComment.UserId,
		UserName:      pbComment.UserName,
		UserAvatarUrl: pbComment.UserAvatarUrl,
		Content:       pbComment.Content,
		Stamp:         pbComment.Stamp,
		Status:        int(pbComment.Status),
		RootId:        pbComment.RootId,
		ReplyId:       pbComment.ReplyId,
	}
	if pbComment.IsRoot {
		comment.IsRoot = 1
	} else {
		comment.IsRoot = 0
	}
}
