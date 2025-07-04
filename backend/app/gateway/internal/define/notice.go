package define

// QueryNoticeListParams 公告列表分页查询参数
type QueryNoticeListParams struct {
	Page     int32  `json:"page" form:"page" binding:"required"`
	PageSize int32  `json:"page_size" form:"page_size" binding:"required"`
	KeyWord  string `json:"keyword" form:"keyword"`
}

type NoticeForm struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}
