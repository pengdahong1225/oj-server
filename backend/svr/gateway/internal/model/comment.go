package model

type CreateCommentRequest struct {
	ObjId         int64  `json:"obj_id"  binding:"required"`
	UserId        int64  `json:"user_id"  binding:"required"`
	UserName      string `json:"user_name"`
	UserAvatarUrl string `json:"user_avatar_url"`
	Content       string `json:"content" binding:"required"`
	Stamp         int64  `json:"stamp"`

	RootId         int64  `json:"root_id"`
	RootCommentId  int64  `json:"root_comment_id"`
	ReplyId        int64  `json:"reply_id"`
	ReplyCommentId int64  `json:"reply_comment_id"`
	ReplyUserName  string `json:"reply_user_name"`
}

// 顶层评论列表查询
type QueryRootCommentListRequest struct {
	ObjId    int64 `form:"obj_id" binding:"required"`
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

// 子评论列表查询
type QueryChildCommentListRequest struct {
	ObjId          int64 `form:"obj_id" binding:"required"`
	RootId         int64 `form:"root_id" binding:"required"`
	RootCommentId  int64 `form:"root_comment_id" binding:"required"`
	ReplyId        int64 `form:"reply_id"`
	ReplyCommentId int64 `form:"reply_comment_id"`
	Cursor         int32 `form:"cursor" binding:"required"`
}

type CommentLikeRequest struct {
	ObjId     int64 `form:"obj_id" binding:"required"`
	CommentId int64 `form:"comment_id" binding:"required"`
}

type DeleteCommentRequest struct {
	ObjId     int64 `form:"obj_id" binding:"required"`
	CommentId int64 `form:"comment_id" binding:"required"`
}
