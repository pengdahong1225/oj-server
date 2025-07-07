package processor

import (
	"oj-server/app/judge/internal/define"
	"oj-server/proto/pb"
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
