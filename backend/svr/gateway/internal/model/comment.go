package model

import "oj-server/pkg/proto/pb"

type Comment struct {
	ID             int64  `json:"id"`
	ObjId          int64  `json:"obj_id"`
	UserId         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	UserAvatarUrl  string `json:"user_avatar_url"`
	Content        string `json:"content"`
	Status         int    `json:"status"`
	ReplyCount     int    `json:"reply_count"`
	LikeCount      int    `json:"like_count"`
	ChildCount     int    `json:"child_count"`
	PubStamp       int64  `json:"pub_stamp"`
	PubRegion      string `json:"pub_region"`
	IsRoot         int    `json:"is_root"`
	RootId         int64  `json:"root_id"`
	RootCommentId  int64  `json:"root_comment_id"`
	ReplyId        int64  `json:"reply_id"`
	ReplyCommentId int64  `json:"reply_comment_id"`
	ReplyUserName  string `json:"reply_user_name"`
}

func (c *Comment) FromPbComment(pbComment *pb.Comment) {
	c.ID = pbComment.Id
	c.ObjId = pbComment.ObjId
	c.UserId = pbComment.UserId
	c.UserName = pbComment.UserName
	c.UserAvatarUrl = pbComment.UserAvatarUrl
	c.Content = pbComment.Content
	c.Status = int(pbComment.Status)
	c.ReplyCount = int(pbComment.ReplyCount)
	c.LikeCount = int(pbComment.LikeCount)
	c.ChildCount = int(pbComment.ChildCount)
	c.PubStamp = pbComment.PubStamp
	c.PubRegion = pbComment.PubRegion
	c.IsRoot = int(pbComment.IsRoot)
	c.RootId = pbComment.RootId
	c.RootCommentId = pbComment.RootCommentId
	c.ReplyId = pbComment.ReplyId
	c.ReplyCommentId = pbComment.ReplyCommentId
	c.ReplyUserName = pbComment.ReplyUserName
}

type CreateCommentForm struct {
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

// 顶层评论列表查询参数
type QueryRootCommentListParams struct {
	ObjId int64 `form:"obj_id" binding:"required"`

	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}
type QueryRootCommentListResult struct {
	Total int64      `json:"total"`
	List  []*Comment `json:"list"`
}

type QueryChildCommentListParams struct {
	ObjId int64 `form:"obj_id" binding:"required"`

	RootId        int64 `form:"root_id" binding:"required"`
	RootCommentId int64 `form:"root_comment_id" binding:"required"`

	ReplyId        int64 `form:"reply_id"`
	ReplyCommentId int64 `form:"reply_comment_id"`

	Cursor int32 `form:"cursor" binding:"required"`
}
type QueryChildCommentListResult struct {
	Total  int64      `json:"total"`
	List   []*Comment `json:"list"`
	Cursor int32      `json:"cursor"`
}

type CommentLikeForm struct {
	ObjId     int64 `form:"obj_id" binding:"required"`
	CommentId int64 `form:"comment_id" binding:"required"`
}

type DeleteCommentForm struct {
	ObjId     int64 `form:"obj_id" binding:"required"`
	CommentId int64 `form:"comment_id" binding:"required"`
}
