package processor

import (
	"oj-server/proto/pb"
	"oj-server/svr/judge/internal/biz"
)

type PyProcessor struct {
	BasicProcessor
}

func (cp *PyProcessor) Compile(param *biz.Param) (*biz.SandBoxApiResponse, error) {
	return nil, nil
}
func (cp *PyProcessor) Run(param *biz.Param) {
}
func (cp *PyProcessor) Judge() []*pb.PBResult {
	return nil
}
