package internal

import (
	"encoding/json"
	"judge-service/internal/judge"
	"judge-service/models"
)

// Handle 代码运行
func Handle(form *models.JudgeRequest) []byte {
	handler := judge.NewHandler()
	rsp := handler.JudgeQuestion(form)
	msg, _ := json.Marshal(rsp)
	return msg
}
