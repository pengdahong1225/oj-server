package processor

import (
	"oj-server/module/proto/pb"
	"oj-server/svr/judge/internal/define"
)

type PyProcessor struct {
	BaseProcessor
}

func (cp *PyProcessor) Compile(param *define.Param) (*define.SandBoxApiResponse, error) {
	return nil, nil
}
func (cp *PyProcessor) Run(param *define.Param) {
}
func (cp *PyProcessor) Judge() []*pb.PBResult {
	return nil
}
