package judge

import (
	"encoding/json"
	"judge-service/models"
)

// Handle 代码运行
func Handle(form *models.JudgeRequest) []byte {
	handler := NewHandler()
	rsp := handler.JudgeQuestion(form)
	msg, _ := json.Marshal(rsp)
	return msg
}
