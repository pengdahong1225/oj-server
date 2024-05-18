package judgeClient

import (
	"google.golang.org/protobuf/proto"
	pb "question-service/logic/proto"
)

var kHeaderLen = 4

func decode(data []byte) (*pb.SSJudgeResponse, error) {
	response := &pb.SSJudgeResponse{}
	err := proto.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func encode(request *pb.SSJudgeRequest) ([]byte, error) {
	ret, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
