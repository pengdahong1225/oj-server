package mysql

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ObjId          int64  `gorm:"column:obj_id" json:"obj_id"`
	UserId         int64  `gorm:"column:user_id" json:"user_id"`
	UserName       string `gorm:"column:user_name" json:"user_name"`
	UserAvatarUrl  string `gorm:"column:user_avatar_url" json:"user_avatar_url"`
	RootId         int64  `gorm:"column:root_id" json:"root_id"`
	RootCommentId  int64  `gorm:"column:root_comment_id" json:"root_comment_id"`
	ReplyId        int64  `gorm:"column:reply_id" json:"reply_id"`
	ReplyCommentId int64  `gorm:"column:reply_comment_id" json:"reply_comment_id"`
	Content        string `gorm:"column:content" json:"content"`
	IsRoot         int    `gorm:"column:is_root" json:"is_root"`
	Status         int    `gorm:"column:status" json:"status"`
	ReplyCount     int    `gorm:"column:reply_count" json:"reply_count"`
	LikeCount      int    `gorm:"column:like_count" json:"like_count"`
	ChildCount     int    `gorm:"column:child_count" json:"child_count"`
	Stamp          int64  `gorm:"column:stamp" json:"stamp"`
}

func (receiver *Comment) TableName() string {
	return "comment"
}
