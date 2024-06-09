package handler

import "errors"

var (
	NotFound      = errors.New("not Found")
	QueryField    = errors.New("query Field")
	AlreadyExists = errors.New("already Exists")
	InsertFailed  = errors.New("insert Failed")
	DeleteFailed  = errors.New("delete Failed")
	UpdateFailed  = errors.New("update Failed")
)
