package comment

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/common/errs"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type CommentServer struct {
	pb.UnimplementedCommentServiceServer
	checker Checker
	saver   Saver
	querier Querier
}

// QueryRootComment
// 查询顶层评论列表，采用偏移量分页
func (c *CommentServer) QueryRootComment(ctx context.Context, in *pb.QueryRootCommentRequest) (*pb.QueryRootCommentResponse, error) {
	// 校验
	if !c.checker.assertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, errs.QueryFailed
	}

	response := &pb.QueryRootCommentResponse{}

	// 查询总量
	db := mysql.DBSession
	var count int64 = 0
	result := db.Model(&mysql.Comment{}).Where("obj_id = ?", in.ObjId).Where("is_root = ?", 1).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	response.Total = int32(count)

	comments := c.querier.onRootComment(in.ObjId, in.Page, in.PageSize)
	for _, comment := range comments {
		response.Data = append(response.Data, translateComment(&comment))
	}
	return response, nil
}

// QueryChildComment
// 查询第二层评论，采用游标分页
func (c *CommentServer) QueryChildComment(ctx context.Context, in *pb.QueryChildCommentRequest) (*pb.QueryChildCommentResponse, error) {
	// 校验
	if !c.checker.assertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, errs.QueryFailed
	}
	if !c.checker.assertRoot(in.RootCommentId, in.RootId) {
		logrus.Errorf("root comment[%d] assert failed", in.RootCommentId)
		return nil, errs.QueryFailed
	}
	if in.ReplyId > 0 && in.ReplyCommentId > 0 && !c.checker.assertReply(in.ReplyCommentId, in.ReplyId) {
		logrus.Errorf("reply comment[%d] assert failed", in.ReplyCommentId)
		return nil, errs.QueryFailed
	}

	response := &pb.QueryChildCommentResponse{}

	// 查询总量
	db := mysql.DBSession
	var count int64 = 0
	result := db.Model(&mysql.Comment{}).Where("obj_id = ?", in.ObjId).Where("is_root = ?", 0).Where("root_id = ?", in.RootId).Where("root_comment_id = ?", in.RootCommentId).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	response.Total = int32(count)

	comments := c.querier.onChildComment(in.ObjId, in.RootId, in.RootCommentId, in.Cursor)
	for _, comment := range comments {
		response.Data = append(response.Data, translateComment(&comment))
	}
	return response, nil
}

func (c *CommentServer) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (c *CommentServer) SaveComment(ctx context.Context, in *pb.SaveCommentRequest) (*emptypb.Empty, error) {
	checker := Checker{}
	comment := in.Data
	if !checker.assertObj(comment.ObjId) {
		logrus.Errorf("obj[%d] assert failed", comment.ObjId)
		return nil, errs.SaveFailed
	}
	if comment.IsRoot > 0 {
		c.saver.onRootComment(comment) // 第一层
	} else {
		c.saver.onChildComment(comment) // 第二层
	}

	return &emptypb.Empty{}, nil
}

func translateComment(comment *mysql.Comment) *pb.Comment {
	return &pb.Comment{
		Id:            comment.ID,
		ObjId:         comment.ObjId,
		UserId:        comment.UserId,
		UserName:      comment.UserName,
		UserAvatarUrl: comment.UserAvatarUrl,
		Content:       comment.Content,
		Status:        int32(comment.Status),
		LikeCount:     int32(comment.LikeCount),
		ReplyCount:    int32(comment.ReplyCount),
		ChildCount:    int32(comment.ChildCount),
		PubStamp:      comment.PubStamp,
		PubRegion:     comment.PubRegion,

		IsRoot:         int32(comment.IsRoot),
		RootId:         comment.RootId,
		RootCommentId:  comment.RootCommentId,
		ReplyId:        comment.ReplyId,
		ReplyCommentId: comment.ReplyCommentId,
	}
}

// CommentLike 评论点赞
func (c *CommentServer) CommentLike(ctx context.Context, in *pb.CommentLikeRequest) (*emptypb.Empty, error) {
	// 校验
	if !c.checker.assertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, errs.QueryFailed
	}

	/*
		update comment set like=like+1 where id=?
	*/
	db := mysql.DBSession
	result := db.Model(&mysql.Comment{}).Where("obj_id = ?", in.ObjId).Where("id = ?", in.CommentId).Update("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.UpdateFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}
	return &empty.Empty{}, nil
}
