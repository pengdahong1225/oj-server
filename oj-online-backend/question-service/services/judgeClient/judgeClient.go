package judgeClient

import (
	"errors"
	"net"
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

func (receiver *TcpClient) Send(msg []byte) (error, []byte) {
	// 检查连接是否可用
	if receiver.conn == nil {
		return errors.New("connection is nil"), nil
	}
	defer receiver.conn.Close() // 用完后关闭连接
	_, err := receiver.conn.Write(msg)
	if err != nil {
		return err, nil
	}
	// 读取服务器返回的数据
	buffer := make([]byte, 2048)
	for {
		n, err := receiver.conn.Read(buffer)
		if err != nil {
			return err, nil
		}
		if n == 0 {
			break
		}
	}

	return nil, buffer
}
