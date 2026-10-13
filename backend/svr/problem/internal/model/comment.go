package model

import (
	"oj-server/common/proto/pb"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreateAt  time.Time      `gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time      `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"`

	ObjID         int64  `gorm:"column:obj_id;not null;index"`
	UserID        int64  `gorm:"column:user_id;not null;index"`
	UserName      string `gorm:"column:user_name;default:''"`
	UserAvatarURL string `gorm:"column:user_avatar_url;default:''"`
	Content       string `gorm:"column:content;type:text;not null"`
	Status        int8   `gorm:"column:status;default:1"`
	ReplyCount    int32  `gorm:"column:reply_count;default:0"`
	LikeCount     int32  `gorm:"column:like_count;default:0"`
	ChildCount    int32  `gorm:"column:child_count;default:0"`
	PubStamp      int64  `gorm:"column:pub_stamp;not null"`
	PubRegion     string `gorm:"column:pub_region;default:''"`

	IsRoot         int8   `gorm:"column:is_root;default:1"`
	RootID         int64  `gorm:"column:root_id;not null"`
	RootCommentID  int64  `gorm:"column:root_comment_id;not null"`
	ReplyID        int64  `gorm:"column:reply_id;not null"`
	ReplyCommentID int64  `gorm:"column:reply_comment_id;not null"`
	ReplyUserName  string `gorm:"column:reply_user_name;default:''"`
}

func (receiver *Comment) TableName() string {
	return "comment"
}

func (receiver *Comment) Transform2Pb() *pb.Comment {
	return &pb.Comment{
		Id:            receiver.ID,
		ObjId:         receiver.ObjID,
		UserId:        receiver.UserID,
		UserName:      receiver.UserName,
		UserAvatarUrl: receiver.UserAvatarURL,
		Content:       receiver.Content,
		Status:        int32(receiver.Status),
		ReplyCount:    int32(receiver.ReplyCount),
		LikeCount:     int32(receiver.LikeCount),
		ChildCount:    int32(receiver.ChildCount),
	}
}
