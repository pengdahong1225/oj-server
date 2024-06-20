package judgeClient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"io"
	"question-service/api/proto"
)

var kHeaderLen = 4

func decode(data []byte) (*pb.SSJudgeResponse, error) {
	header := make([]byte, kHeaderLen)
	_, err := io.ReadFull(bytes.NewReader(data), header)
	if err != nil {
		logrus.Errorf("read header error: %s", err.Error())
		return nil, err
	}
	// 解码包头，得到数据长度
	length := binary.BigEndian.Uint32(header)
	if length != uint32(len(data)-kHeaderLen) {
		logrus.Errorf("the length of package is error , length: %d", length)
		return nil, errors.New("the length of package is error")
	}

	// 解码数据
	body := make([]byte, length)
	_, err = io.ReadFull(bytes.NewReader(data), body)
	if err != nil {
		logrus.Errorf("the body of package is error , err: %s", err.Error())
		return nil, err
	}

	// 协议解析
	response := &pb.SSJudgeResponse{}
	err = proto.Unmarshal(body, response)
	if err != nil {
		logrus.Errorf("the body of package is error , err: %s", err.Error())
		return nil, err
	}
	return response, nil
}

func encode(request *pb.SSJudgeRequest) ([]byte, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		logrus.Errorf("marshal request failed, err: %s", err.Error())
		return nil, err
	}

	length := uint32(len(data))
	buf := new(bytes.Buffer)
	// 将数据长度编码为4字节的大端整数
	binary.Write(buf, binary.BigEndian, length)

	// 写入body数据
	buf.Write(data)

	return buf.Bytes(), nil
}
