package processor

import (
	"oj-server/proto/pb"
	"oj-server/svr/judge/internal/define"
)

type GoProcessor struct {
	BaseProcessor
}

func (cp *GoProcessor) Compile(param *define.Param) (*define.SandBoxApiResponse, error) {
	return nil, nil
}
func (cp *GoProcessor) Run(param *define.Param) {
}
func (cp *GoProcessor) Judge() []*pb.PBResult {
	return nil
}
