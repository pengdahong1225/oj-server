package model

type AddNoticeRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoticeRequest struct {
	ID      int64  `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type GetNoticeListRequest struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}
