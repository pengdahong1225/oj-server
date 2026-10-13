package service

import (
	"context"
	"oj-server/common/proto/pb"
	"oj-server/pkg/utils"
	"oj-server/svr/problem/internal/client"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// comment服务
type CommentService struct {
	pb.UnimplementedCommentServiceServer
	repo    *data.CommentRepo
	clients *client.Clients
}

func NewCommentService(clients *client.Clients) (*CommentService, error) {
	repo, err := data.NewCommentRepo()
	if err != nil {
		logrus.Errorf("NewProblemService failed, err:%s", err.Error())
		return nil, err
	}

	return &CommentService{
		repo:    repo,
		clients: clients,
	}, nil
}

func (cs *CommentService) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*emptypb.Empty, error) {
	// 1.构造评论
	comment := &model.Comment{
		ObjID:         in.ObjId,
		UserID:        in.UserId,
		UserName:      in.UserName,
		UserAvatarURL: in.UserAvatarUrl,
		Content:       in.Content,
		PubStamp:      in.Stamp,
	}
	// 非楼主评论
	if in.RootId > 0 && in.RootCommentId > 0 {
		comment.IsRoot = 0
		comment.RootID = in.RootId
		comment.RootCommentID = in.RootCommentId
		if in.ReplyId > 0 && in.ReplyCommentId > 0 {
			comment.ReplyID = in.ReplyId
			comment.ReplyCommentID = in.ReplyCommentId
			comment.ReplyUserName = in.ReplyUserName
		} else {
			// 默认回复楼主
			comment.ReplyID = in.RootId
			comment.ReplyCommentID = in.RootCommentId
			comment.ReplyUserName = in.ReplyUserName
		}
	} else {
		// 默认楼主评论
		comment.RootID = 0
		comment.RootCommentID = 0
		comment.IsRoot = 1
		comment.ReplyID = 0
		comment.ReplyCommentID = 0
	}
	// todo 查询ip归属地会影响响应时间，改为异步
	go func() {
		// ip归属地
		info, err := utils.QueryIpGeolocation(in.Address)
		if err != nil {
			logrus.Errorf("查询ip归属地失败,ip:%s,err:%s", in.Address, err.Error())
		}
		comment.PubRegion = info.RegionName

		// 2.校验评论区
		if !cs.repo.AssertObj(in.ObjId) {
			logrus.Errorf("obj[%d] assert failed", in.ObjId)
			return
		}
		// 3.保存评论
		if comment.IsRoot > 0 {
			cs.repo.SaveRootComment(comment) // 第一层
		} else {
			cs.repo.SaveChildComment(comment) // 第二层
		}
	}()

	return nil, nil
}

func (cs *CommentService) QueryRootComment(ctx context.Context, in *pb.QueryRootCommentRequest) (*pb.QueryRootCommentResponse, error) {
	// 校验
	if !cs.repo.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := cs.repo.QueryRootComment(in.ObjId, int(in.Page), int(in.PageSize))
	if err != nil {
		logrus.Errorf("query root comment failed, err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "query root comment failed")
	}

	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:             comment.ID,
			ObjId:          comment.ObjID,
			UserId:         comment.UserID,
			UserName:       comment.UserName,
			UserAvatarUrl:  comment.UserAvatarURL,
			Content:        comment.Content,
			ChildCount:     int32(comment.ChildCount),
			ReplyCount:     int32(comment.ReplyCount),
			LikeCount:      int32(comment.LikeCount),
			PubStamp:       comment.PubStamp,
			PubRegion:      comment.PubRegion,
			IsRoot:         int32(comment.IsRoot),
			RootId:         comment.RootID,
			RootCommentId:  comment.RootCommentID,
			ReplyId:        comment.ReplyID,
			ReplyCommentId: comment.ReplyCommentID,
			ReplyUserName:  comment.ReplyUserName,
		})
	}

	return &pb.QueryRootCommentResponse{
		Total: int32(total),
		Data:  pbComments,
	}, nil
}
func (cs *CommentService) QueryChildComment(ctx context.Context, in *pb.QueryChildCommentRequest) (*pb.QueryChildCommentResponse, error) {
	// 校验
	if !cs.repo.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := cs.repo.QueryChildComment(in.ObjId, in.RootId, in.RootCommentId, in.Cursor)
	if err != nil {
		logrus.Errorf("query child comment failed, err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "query child comment failed")
	}

	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:             comment.ID,
			ObjId:          comment.ObjID,
			UserId:         comment.UserID,
			UserName:       comment.UserName,
			UserAvatarUrl:  comment.UserAvatarURL,
			Content:        comment.Content,
			ReplyCount:     int32(comment.ReplyCount),
			LikeCount:      int32(comment.LikeCount),
			PubStamp:       comment.PubStamp,
			PubRegion:      comment.PubRegion,
			IsRoot:         int32(comment.IsRoot),
			RootId:         comment.RootID,
			RootCommentId:  comment.RootCommentID,
			ReplyId:        comment.ReplyID,
			ReplyCommentId: comment.ReplyCommentID,
			ReplyUserName:  comment.ReplyUserName,
		})
	}
	return &pb.QueryChildCommentResponse{
		Total: int32(total),
		Data:  pbComments,
	}, nil
}
func (cs *CommentService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (cs *CommentService) CommentLike(ctx context.Context, in *pb.CommentLikeRequest) (*emptypb.Empty, error) {
	// 校验
	if !cs.repo.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	cs.repo.CommentLike(in.ObjId, in.CommentId)
	return &emptypb.Empty{}, nil
}
