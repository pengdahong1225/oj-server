package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
	"oj-server/utils"
)

// comment服务
type CommentService struct {
	pb.UnimplementedCommentServiceServer
	uc *biz.CommentUseCase
}

func NewCommentService() *CommentService {
	repo, err := data.NewCommentRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	return &CommentService{
		uc: biz.NewCommentUseCase(repo), // 注入实现
	}
}

func (cs *CommentService) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*emptypb.Empty, error) {
	// 1.构造评论
	comment := &model.Comment{
		ObjId:         in.ObjId,
		UserId:        in.UserId,
		UserName:      in.UserName,
		UserAvatarUrl: in.UserAvatarUrl,
		Content:       in.Content,
		PubStamp:      in.Stamp,
	}
	// 非楼主评论
	if in.RootId > 0 && in.RootCommentId > 0 {
		comment.IsRoot = 0
		comment.RootId = in.RootId
		comment.RootCommentId = in.RootCommentId
		if in.ReplyId > 0 && in.ReplyCommentId > 0 {
			comment.ReplyId = in.ReplyId
			comment.ReplyCommentId = in.ReplyCommentId
			comment.ReplyUserName = in.ReplyUserName
		} else {
			// 默认回复楼主
			comment.ReplyId = in.RootId
			comment.ReplyCommentId = in.RootCommentId
			comment.ReplyUserName = in.ReplyUserName
		}
	} else {
		// 默认楼主评论
		comment.RootId = 0
		comment.RootCommentId = 0
		comment.IsRoot = 1
		comment.ReplyId = 0
		comment.ReplyCommentId = 0
	}
	// ip归属地
	ip := ctx.Value("address")
	info, err := utils.QueryIpGeolocation(ip.(string))
	if err != nil {
		logrus.Errorf("查询ip归属地失败,ip:%s,err:%s", ip, err.Error())
	}
	comment.PubRegion = info.RegionName

	// 2.校验评论区
	if !cs.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}
	// 3.保存评论
	if comment.IsRoot > 0 {
		cs.uc.SaveRootComment(comment) // 第一层
	} else {
		cs.uc.SaveChildComment(comment) // 第二层
	}

	return nil, nil
}

func (cs *CommentService) QueryRootComment(ctx context.Context, in *pb.QueryRootCommentRequest) (*pb.QueryRootCommentResponse, error) {
	// 校验
	if !cs.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := cs.uc.QueryRootComment(in.ObjId, int(in.Page), int(in.PageSize))
	if err != nil {
		logrus.Errorf("query root comment failed, err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "query root comment failed")
	}

	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:             comment.ID,
			ObjId:          comment.ObjId,
			UserId:         comment.UserId,
			UserName:       comment.UserName,
			UserAvatarUrl:  comment.UserAvatarUrl,
			Content:        comment.Content,
			ReplyCount:     int32(comment.ReplyCount),
			LikeCount:      int32(comment.LikeCount),
			PubStamp:       comment.PubStamp,
			PubRegion:      comment.PubRegion,
			IsRoot:         int32(comment.IsRoot),
			RootId:         comment.RootId,
			RootCommentId:  comment.RootCommentId,
			ReplyId:        comment.ReplyId,
			ReplyCommentId: comment.ReplyCommentId,
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
	if !cs.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := cs.uc.QueryChildComment(in.ObjId, in.RootId, in.RootCommentId, in.Cursor)
	if err != nil {
		logrus.Errorf("query child comment failed, err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "query child comment failed")
	}

	pbComments := make([]*pb.Comment, 0, len(comments))
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:             comment.ID,
			ObjId:          comment.ObjId,
			UserId:         comment.UserId,
			UserName:       comment.UserName,
			UserAvatarUrl:  comment.UserAvatarUrl,
			Content:        comment.Content,
			ReplyCount:     int32(comment.ReplyCount),
			LikeCount:      int32(comment.LikeCount),
			PubStamp:       comment.PubStamp,
			PubRegion:      comment.PubRegion,
			IsRoot:         int32(comment.IsRoot),
			RootId:         comment.RootId,
			RootCommentId:  comment.RootCommentId,
			ReplyId:        comment.ReplyId,
			ReplyCommentId: comment.ReplyCommentId,
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
	if !cs.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	cs.uc.CommentLike(in.ObjId, in.CommentId)
	return &emptypb.Empty{}, nil
}
