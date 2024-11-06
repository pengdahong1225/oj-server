package types

import "github.com/pengdahong1225/oj-server/backend/proto/pb"

type Param struct {
	Uid           int64
	ProblemID     int64
	Code          string // 源代码
	Language      string // 语种
	ProblemConfig *pb.ProblemConfig

	// 编译结果中读取
	FileIds map[string]string // 文件id

	Accepted bool
	Message  string
}
