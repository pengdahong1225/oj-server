package utils

import (
	"github.com/google/uuid"
	"net"
)

// GetOutboundIP 获取本机的出口IP
// UDP 是一种无连接的协议，因此不需要关心连接是否成功到达目的地，一旦尝试发起连接，操作系统将为本次连接设置ip地址。
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53") // 使用 udp 不需要关注是否送达
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
func GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
