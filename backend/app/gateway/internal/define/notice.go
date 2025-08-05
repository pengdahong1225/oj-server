package define

import "oj-server/proto/pb"

// ============================ notice 列表 ============================
type QueryNoticeListParams struct {
	Page     int32  `json:"page" form:"page" binding:"required"`
	PageSize int32  `json:"page_size" form:"page_size" binding:"required"`
	KeyWord  string `json:"keyword" form:"keyword"`
}
type QueryNoticeListResponse struct {
	Total      int32        `json:"total"`
	NoticeList []*pb.Notice `json:"notice_list"`
}

// ============================ notice ============================
type NoticeForm struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}
