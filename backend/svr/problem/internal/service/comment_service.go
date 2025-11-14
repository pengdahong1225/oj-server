package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"oj-server/global"
	"oj-server/pkg/gPool"
	"oj-server/pkg/mq"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
)

// comment服务
type CommentService struct {
	pb.UnimplementedCommentServiceServer
	uc *biz.CommentUseCase

	comment_consumer *mq.Consumer // 评论任务消费者
}

func NewCommentService() *CommentService {
	repo, err := data.NewCommentRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	return &CommentService{
		uc: biz.NewCommentUseCase(repo), // 注入实现
		comment_consumer: mq.NewConsumer(
			global.RabbitMqExchangeKind,
			global.RabbitMqExchangeName,
			global.RabbitMqCommentQueue,
			global.RabbitMqCommentKey,
			"", // 消费者标签，用于区别不同的消费者
		),
	}
}

func (ps *CommentService) ConsumeComment() {
	deliveries := ps.comment_consumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer ps.comment_consumer.Close()

	for d := range deliveries {
		// 处理comment
		result := func(data []byte) bool {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("consume comment panic: %v", err)
				}
			}()
			comment := &pb.Comment{}
			err := proto.Unmarshal(data, comment)
			if err != nil {
				logrus.Errorln("解析comment err：", err.Error())
				return false
			}
			// 异步处理
			_ = gPool.Instance().Submit(func() {
				ps.HandleSaveComment(comment)
			})
			return true
		}(d.Body)

		// 确认
		if result {
			d.Ack(false)
		} else {
			d.Reject(false)
		}
	}
}

func (ps *CommentService) HandleSaveComment(pbComment *pb.Comment) {
	if !ps.uc.AssertObj(pbComment.ObjId) {
		logrus.Errorf("obj[%d] assert failed", pbComment.ObjId)
		return
	}
	// 评论基本信息
	comment := &model.Comment{
		ObjId:         pbComment.ObjId,
		UserId:        pbComment.UserId,
		UserName:      pbComment.UserName,
		UserAvatarUrl: pbComment.UserAvatarUrl,
		Content:       pbComment.Content,
		PubStamp:      pbComment.PubStamp,
		PubRegion:     pbComment.PubRegion,
	}
	// 处理评论的级别信息
	if pbComment.IsRoot > 0 {
		comment.IsRoot = 1
		comment.RootId = pbComment.RootId
		comment.RootCommentId = pbComment.RootCommentId
	} else {
		comment.IsRoot = 0
		if pbComment.ReplyId > 0 && pbComment.ReplyCommentId > 0 {
			comment.ReplyId = pbComment.ReplyId
			comment.ReplyCommentId = pbComment.ReplyCommentId
			comment.ReplyUserName = pbComment.ReplyUserName
		}
	}

	if comment.IsRoot > 0 {
		ps.uc.SaveRootComment(comment) // 第一层
	} else {
		ps.uc.SaveChildComment(comment) // 第二层
	}
}

//// 非楼主评论
//if form.RootId > 0 && form.RootCommentId > 0 {
//	pbComment.IsRoot = 0
//	pbComment.RootId = form.RootId
//	pbComment.RootCommentId = form.RootCommentId
//	if form.ReplyId > 0 && form.ReplyCommentId > 0 {
//		pbComment.ReplyId = form.ReplyId
//		pbComment.ReplyCommentId = form.ReplyCommentId
//		pbComment.ReplyUserName = form.ReplyUserName
//	} else {
//		// 默认回复楼主
//		pbComment.ReplyId = form.RootId
//		pbComment.ReplyCommentId = form.RootCommentId
//		pbComment.ReplyUserName = form.ReplyUserName
//	}
//} else {
//	// 楼主评论
//	pbComment.RootId = 0
//	pbComment.RootCommentId = 0
//	pbComment.IsRoot = 1
//	pbComment.ReplyId = 0
//	pbComment.ReplyCommentId = 0
//}
//// 查询ip归属地
//ip := ctx.ClientIP()
//info, err := utils.QueryIpGeolocation(ip)
//if err != nil {
//	logrus.Errorf("查询ip归属地失败,ip:%s,err:%s", ip, err.Error())
//	info = &utils.IPInfoResp{
//		RegionName: "未知地区",
//	}
//}
//pbComment.PubRegion = info.RegionName

func (ps *CommentService) QueryRootComment(ctx context.Context, in *pb.QueryRootCommentRequest) (*pb.QueryRootCommentResponse, error) {
	// 校验
	if !ps.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := ps.uc.QueryRootComment(in.ObjId, int(in.Page), int(in.PageSize))
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
func (ps *CommentService) QueryChildComment(ctx context.Context, in *pb.QueryChildCommentRequest) (*pb.QueryChildCommentResponse, error) {
	// 校验
	if !ps.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := ps.uc.QueryChildComment(in.ObjId, in.RootId, in.RootCommentId, in.Cursor)
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
func (ps *CommentService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (ps *CommentService) CommentLike(ctx context.Context, in *pb.CommentLikeRequest) (*emptypb.Empty, error) {
	// 校验
	if !ps.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	ps.uc.CommentLike(in.ObjId, in.CommentId)
	return &emptypb.Empty{}, nil
}
