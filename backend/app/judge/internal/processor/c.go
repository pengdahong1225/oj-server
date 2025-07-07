package processor

import (
	"oj-server/app/judge/internal/define"
	"oj-server/proto/pb"
)

type CProcessor struct {
	BaseProcessor
}

func (cp *CProcessor) Compile(param *define.Param) (*define.SandBoxApiResponse, error) {
	return nil, nil
}
func (cp *CProcessor) Run(param *define.Param) {
}
func (cp *CProcessor) Judge() []*pb.PBResult {
	return nil
}
