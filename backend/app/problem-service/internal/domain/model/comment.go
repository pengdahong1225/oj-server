package model

import (
	"time"
)

type Comment struct {
	ID            int64     `gorm:"primary_key"`
	CreateAt      time.Time `gorm:"column:create_at" json:"create_at"`
	UpdateAt      time.Time `gorm:"column:update_at" json:"update_at"`
	DeletedAt     time.Time `gorm:"column:delete_at" json:"delete_at"`
	ObjId         int64     `gorm:"column:obj_id" json:"obj_id"`
	UserId        int64     `gorm:"column:user_id" json:"user_id"`
	UserName      string    `gorm:"column:user_name" json:"user_name"`
	UserAvatarUrl string    `gorm:"column:user_avatar_url" json:"user_avatar_url"`
	Content       string    `gorm:"column:content" json:"content"`
	Status        int       `gorm:"column:status" json:"status"`
	ReplyCount    int       `gorm:"column:reply_count" json:"reply_count"`
	LikeCount     int       `gorm:"column:like_count" json:"like_count"`
	ChildCount    int       `gorm:"column:child_count" json:"child_count"`
	PubStamp      int64     `gorm:"column:pub_stamp" json:"pub_stamp"`
	PubRegion     string    `gorm:"column:pub_region" json:"pub_region"`

	IsRoot         int    `gorm:"column:is_root" json:"is_root"`
	RootId         int64  `gorm:"column:root_id" json:"root_id"`
	RootCommentId  int64  `gorm:"column:root_comment_id" json:"root_comment_id"`
	ReplyId        int64  `gorm:"column:reply_id" json:"reply_id"`
	ReplyCommentId int64  `gorm:"column:reply_comment_id" json:"reply_comment_id"`
	ReplyUserName  string `gorm:"column:reply_user_name" json:"reply_user_name"`
}

func (receiver *Comment) TableName() string {
	return "comment"
}

func NewComment() *Comment {
	return &Comment{
		Status:     1,
		ReplyCount: 0,
		LikeCount:  0,
		ChildCount: 0,
		PubStamp:   time.Now().Unix(),
	}
}
