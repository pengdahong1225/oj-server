package service

import (
	"context"
	"oj-server/common/proto/pb"
	"oj-server/svr/problem/internal/client"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 公告服务
type NoticeService struct {
	pb.UnimplementedNoticeServiceServer
	repo    *data.NoticeRepo
	clients *client.Clients
}

func NewNoticeService(clients *client.Clients) (*NoticeService, error) {
	repo, err := data.NewNoticeRepo()
	if err != nil {
		logrus.Errorf("NewNoticeService failed, err:%s", err.Error())
		return nil, err
	}

	return &NoticeService{
		repo:    repo,
		clients: clients,
	}, nil
}

// CreateNotice 新建公告
func (ns *NoticeService) CreateNotice(ctx context.Context, in *pb.CreateNoticeRequest) (*emptypb.Empty, error) {
	if in.Notice == nil {
		return nil, status.Errorf(codes.InvalidArgument, "notice is nil")
	}
	notice := &model.Notice{
		Title:    in.Notice.Title,
		Content:  in.Notice.Content,
		CreateBy: in.Notice.CreateBy,
		// status 由 gorm 默认值 1（正常/可见）填充
	}
	if err := ns.repo.CreateNotice(notice); err != nil {
		logrus.Errorf("create notice failed, err:%s", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteNotice 删除公告
func (ns *NoticeService) DeleteNotice(ctx context.Context, in *pb.DeleteNoticeRequest) (*emptypb.Empty, error) {
	if err := ns.repo.DeleteNotice(in.Id); err != nil {
		logrus.Errorf("delete notice failed, err:%s", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateNotice 更新公告
func (ns *NoticeService) UpdateNotice(ctx context.Context, in *pb.UpdateNoticeRequest) (*emptypb.Empty, error) {
	if in.Notice == nil {
		return nil, status.Errorf(codes.InvalidArgument, "notice is nil")
	}
	notice := &model.Notice{
		ID:      in.Notice.Id,
		Title:   in.Notice.Title,
		Content: in.Notice.Content,
	}
	if err := ns.repo.UpdateNotice(notice); err != nil {
		logrus.Errorf("update notice failed, err:%s", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// PublishNotice 发布公告
func (ns *NoticeService) PublishNotice(ctx context.Context, in *pb.PublishNoticeRequest) (*emptypb.Empty, error) {
	if err := ns.repo.UpdateNoticeStatus(in.Id, 1); err != nil {
		logrus.Errorf("publish notice failed, err:%s", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// HideNotice 隐藏公告
func (ns *NoticeService) HideNotice(ctx context.Context, in *pb.HideNoticeRequest) (*emptypb.Empty, error) {
	if err := ns.repo.UpdateNoticeStatus(in.Id, 0); err != nil {
		logrus.Errorf("hide notice failed, err:%s", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetNoticeList 分页查询公告列表
func (ns *NoticeService) GetNoticeList(ctx context.Context, in *pb.GetNoticeListRequest) (*pb.GetNoticeListResponse, error) {
	total, notices, err := ns.repo.QueryNoticeList(int(in.Page), int(in.PageSize))
	if err != nil {
		logrus.Errorf("query notice list failed, err:%s", err.Error())
		return nil, err
	}

	pbNotices := make([]*pb.Notice, 0, len(notices))
	for i := range notices {
		pbNotices = append(pbNotices, notices[i].Transform2Pb())
	}

	return &pb.GetNoticeListResponse{
		Total: total,
		List:  pbNotices,
	}, nil
}

// GetNotice 查询公告详情
func (ns *NoticeService) GetNotice(ctx context.Context, in *pb.GetNoticeRequest) (*pb.GetNoticeResponse, error) {
	notice, err := ns.repo.QueryNotice(in.Id)
	if err != nil {
		logrus.Errorf("query notice failed, err:%s", err.Error())
		return nil, err
	}
	return &pb.GetNoticeResponse{
		Notice: notice.Transform2Pb(),
	}, nil
}
