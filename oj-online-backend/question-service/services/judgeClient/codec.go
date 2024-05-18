package judgeClient

import (
	"google.golang.org/protobuf/proto"
	pb "question-service/logic/proto"
)

var kHeaderLen = 4

func decode(data []byte) *pb.SSJudgeResponse {
	response := &pb.SSJudgeResponse{}
	err := proto.Unmarshal(data, response)
	if err != nil {
		return nil
	}
	return response
}

func encode(request *pb.SSJudgeRequest) []byte {
	ret, err := proto.Marshal(request)
	if err != nil {
		return nil
	}
	return ret
}
