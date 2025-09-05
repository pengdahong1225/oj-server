package model

// binding的逗号之间不能有空格

type AddCommentForm struct {
	ObjId         int64  `json:"obj_id" form:"obj_id" binding:"required"`
	UserId        int64  `json:"user_id" form:"user_id" binding:"required"`
	UserName      string `json:"user_name" form:"user_name"`
	UserAvatarUrl string `json:"user_avatar_url" form:"user_avatar_url"`
	Content       string `json:"content" form:"content" binding:"required"`
	Stamp         int64  `json:"stamp" form:"stamp"`

	RootId         int64  `json:"root_id" form:"root_id"`
	RootCommentId  int64  `json:"root_comment_id" form:"root_comment_id"`
	ReplyId        int64  `json:"reply_id" form:"reply_id"`
	ReplyCommentId int64  `json:"reply_comment_id" form:"reply_comment_id"`
	ReplyUserName  string `json:"reply_user_name" form:"reply_user_name"`
}

// RootCommentListQueryParams 顶层评论列表查询参数
type RootCommentListQueryParams struct {
	ObjId int64 `json:"obj_id" form:"obj_id" binding:"required"`

	Page     int32 `json:"page" form:"page" binding:"required"`
	PageSize int32 `json:"page_size" form:"page_size" binding:"required"`
}

type ChildCommentListQueryParams struct {
	ObjId int64 `json:"obj_id" form:"obj_id" binding:"required"`

	RootId        int64 `json:"root_id" form:"root_id" binding:"required"`
	RootCommentId int64 `json:"root_comment_id" form:"root_comment_id" binding:"required"`

	ReplyId        int64 `json:"reply_id" form:"reply_id"`
	ReplyCommentId int64 `json:"reply_comment_id" form:"reply_comment_id"`

	Cursor int32 `json:"cursor" form:"cursor" binding:"required"`
}

type CommentLikeForm struct {
	ObjId     int64 `json:"obj_id" form:"obj_id" binding:"required"`
	CommentId int64 `json:"comment_id" form:"comment_id" binding:"required"`
}
