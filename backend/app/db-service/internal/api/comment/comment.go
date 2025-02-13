package comment

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/common/errs"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentServer struct {
	pb.UnimplementedCommentServiceServer
	checker Checker
	saver   Saver
	querier Querier
}

func (c *CommentServer) QueryComment(ctx context.Context, in *pb.QueryCommentRequest) (*pb.QueryCommentResponse, error) {
	// 校验
	if !c.checker.assertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, errs.QueryFailed
	}
	if in.RootId > 0 && in.RootCommentId > 0 && !c.checker.assertRoot(in.RootCommentId, in.RootId) {
		logrus.Errorf("root comment[%d] assert failed", in.RootCommentId)
		return nil, errs.QueryFailed
	}
	if in.ReplyId > 0 && in.ReplyCommentId > 0 && !c.checker.assertReply(in.ReplyCommentId, in.ReplyId) {
		logrus.Errorf("reply comment[%d] assert failed", in.ReplyCommentId)
		return nil, errs.QueryFailed
	}

	response := &pb.QueryCommentResponse{}
	cursor := in.Cursor
	if cursor < 0 {
		cursor = 0
	}
	if in.RootId == 0 || in.RootCommentId == 0 {
		comments := c.querier.onRootComment(in.ObjId, in.Cursor)
		for _, comment := range comments {
			response.Data = append(response.Data, translateComment(&comment))
		}
		return response, nil
	} else {
		comments := c.querier.onChildComment(in.ObjId, in.Cursor)
		for _, comment := range comments {
			response.Data = append(response.Data, translateComment(&comment))
		}
		return response, nil
	}
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
