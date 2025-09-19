package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"oj-server/global"
	"oj-server/module/gPool"
	"oj-server/module/mq"
	"oj-server/module/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
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
			comment := &pb.Comment{}
			err := proto.Unmarshal(data, comment)
			if err != nil {
				logrus.Errorln("解析comment err：", err.Error())
				return false
			}
			// 异步处理
			_ = gPool.Instance().Submit(func() {
				ps.HandleAddComment(comment)
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

func (ps *CommentService) HandleAddComment(comment *pb.Comment) {
	if !ps.uc.AssertObj(comment.ObjId) {
		logrus.Errorf("obj[%d] assert failed", comment.ObjId)
		return
	}
	if comment.IsRoot > 0 {
		ps.uc.SaveRootComment(comment) // 第一层
	} else {
		ps.uc.SaveChildComment(comment) // 第二层
	}
}

func (ps *CommentService) QueryRootComment(ctx context.Context, in *pb.QueryRootCommentRequest) (*pb.QueryRootCommentResponse, error) {
	// 校验
	if !ps.uc.AssertObj(in.ObjId) {
		logrus.Errorf("obj[%d] assert failed", in.ObjId)
		return nil, status.Errorf(codes.NotFound, "obj[%d] assert failed", in.ObjId)
	}

	total, comments, err := ps.uc.QueryRootComment(int(in.Page), int(in.PageSize))
	if err != nil {
		logrus.Errorf("query root comment failed, err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "query root comment failed")
	}
	return &pb.QueryRootCommentResponse{
		Total: total,
		Data:  comments,
	}, nil
}
func (ps *CommentService) QueryChildComment(ctx context.Context, in *pb.QueryChildCommentRequest) (*pb.QueryChildCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryChildComment not implemented")
}
func (ps *CommentService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (ps *CommentService) SaveComment(ctx context.Context, in *pb.SaveCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveComment not implemented")
}
func (ps *CommentService) CommentLike(ctx context.Context, in *pb.CommentLikeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentLike not implemented")
}
