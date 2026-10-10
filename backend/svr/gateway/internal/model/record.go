package model

type QueryRecordListRequest struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}
