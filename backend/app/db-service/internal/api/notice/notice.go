package notice

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/common/errs"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"strconv"
)

type NoticeServer struct {
	pb.UnimplementedNoticeServiceServer
}

func (r *NoticeServer) GetNoticeList(ctx context.Context, request *pb.GetNoticeListRequest) (*pb.GetNoticeListResponse, error) {
	db := mysql.DBSession
	rsp := &pb.GetNoticeListResponse{}
	name := "%" + request.Keyword + "%"
	offSet := int((request.Page - 1) * request.PageSize)

	/*
		select COUNT(*) AS count from notice
		where title like '%name%';
	*/
	var count int64 = 0
	result := db.Model(&mysql.Notice{}).Where("title LIKE ?", name).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	rsp.Total = int32(count)
	logrus.Debugln("count:", count)

	/*
		select id,title,content,create_at,status,create_by from notice
		where title like '%name%'
		order by id
		offset off_set
		limit page_size;
	*/
	var noticeList []mysql.Notice
	result = db.Select("id,title,content,create_at,status, create_by").Where("title LIKE ?", name).Order("id").Offset(offSet).Limit(int(request.PageSize)).Find(&noticeList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	for _, v := range noticeList {
		rsp.Data = append(rsp.Data, &pb.Notice{
			Id:       v.ID,
			Title:    v.Title,
			CreateAt: strconv.FormatInt(v.CreateAt.Unix(), 10),
			Content:  v.Content,
			Status:   v.Status,
			CreateBy: v.CreateBy,
		})
	}
	return rsp, nil
}
