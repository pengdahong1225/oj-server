package judgeClient

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"question-service/api/proto"
	"time"
)

// TcpClient 负责包的发送和接受
type TcpClient struct {
	conn *net.TCPConn
}

func (receiver *TcpClient) Connect(dsn string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", dsn)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	receiver.conn = conn
	return nil
}

func (receiver *TcpClient) Request(req *pb.SSJudgeRequest) (*pb.SSJudgeResponse, error) {
	// 检查连接是否可用
	if receiver.conn == nil {
		return nil, errors.New("connection is nil")
	}
	defer receiver.conn.Close() // 用完后关闭连接

	// 编码
	msg, err := encode(req)
	if err != nil {
		return nil, err
	}
	// 发送
	_, err = receiver.conn.Write(msg)
	if err != nil {
		return nil, err
	}

	// 读取（带超时）
	rspChan := make(chan *pb.SSJudgeResponse)
	errChan := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置5秒超时
	defer cancel()
	go func(ctx context.Context) {
		// 读取
		buffer := make([]byte, 2048)
		n, err := receiver.conn.Read(buffer)
		if err != nil {
			logrus.Infoln("read error:", err)
			errChan <- err
		} else {
			rep, err := decode(buffer[:n]) // 操作读取的数据必须带上实际n，读多少用多少
			if err != nil {
				logrus.Infoln("decode error:", err)
				errChan <- err
			} else {
				rspChan <- rep
			}
		}
	}(ctx)

	// 阻塞等待
	select {
	case rsp := <-rspChan:
		return rsp, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
