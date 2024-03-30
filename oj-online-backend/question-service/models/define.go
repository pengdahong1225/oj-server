package models

// JudgeRequest question-service -->> MQ -->> judge-service
type JudgeRequest struct {
	SessionID  string `json:"sessionID" form:"sessionID"`
	QuestionID int64  `json:"questionID" form:"questionID"`
	UserID     int64  `json:"userID" form:"userID"`
	Title      string `json:"title" form:"title"`
	Code       string `json:"code" form:"code"`
	Clang      string `json:"clang" form:"clang"`
}
