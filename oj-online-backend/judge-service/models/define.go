package models

const (
	EN_Status_OK           = iota // 正确
	EN_Status_Faild               // 答案错误
	EN_Status_MemoryOver          // 内存溢出
	EN_Status_CompileError        // 编译错误
	EN_Status_TimeOut             // 时间超时
	EN_Status_Internal            // 系统错误
)

type TestCase struct {
	Input  string `json:"input"`  // 输入
	Output string `json:"output"` // 输出
}

var ValidGolangPackageMap = map[string]struct{}{
	"bytes":   {},
	"fmt":     {},
	"math":    {},
	"sort":    {},
	"strings": {},
}

type JudgeBack struct {
	SessionID  string `json:"sessionID" form:"id"`
	QuestionID int64  `json:"questionID" form:"questionID"`
	UserID     int64  `json:"userID" form:"userID"`
	Clang      string `json:"clang" form:"clang"`
	Status     int32  `json:"status" form:"status"` // 0: 正常 1: 代码非法 2: 编译错误 3: 运行超时 4: 内存溢出 5: 系统错误
	Tips       string `json:"tips" form:"tips"`     // 系统提示
	Output     string `json:"output" form:"output"` // 系统输出
}
