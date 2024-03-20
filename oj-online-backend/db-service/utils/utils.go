package utils

import (
	"net"
	"strings"
)

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53") // 使用 udp 不需要关注是否送达
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// SplitStringWithX 分割字符串
func SplitStringWithX(src string, X string) []string {
	// 去除头尾空格
	str := strings.TrimSpace(src)
	// 按'#'分割src #xxx
	result := strings.Split(str, X)
	// 去除result中的空字符串
	for i := 0; i < len(result); i++ {
		if result[i] == "" {
			result = append(result[:i], result[i+1:]...)
			i = 0
		}
	}
	return result
}

// SpliceStringWithX 合并字符串，用X做前缀
func SpliceStringWithX(src []string, X string) string {
	builder := strings.Builder{}
	for _, s := range src {
		builder.WriteString(X)
		builder.WriteString(s)
	}
	return builder.String()
}
