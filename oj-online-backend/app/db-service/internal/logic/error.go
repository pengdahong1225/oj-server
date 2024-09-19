package logic

import "errors"

var (
	NotFound      = errors.New("not Found")
	QueryFailed   = errors.New("query Failed")
	AlreadyExists = errors.New("already Exists")
	InsertFailed  = errors.New("insert Failed")
	DeleteFailed  = errors.New("delete Failed")
	UpdateFailed  = errors.New("update Failed")
)
