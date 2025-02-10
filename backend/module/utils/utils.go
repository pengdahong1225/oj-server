package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
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
func GetOutboundIPString() (string, error) {
	ip, err := GetOutboundIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}
func GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
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

// GenerateSmsCode 生成长度指定长度的验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// GenerateSubmitID 生成提交ID
func GenerateSubmitID(userID, questionID int) (int64, error) {
	// 将用户ID和问题ID转换为字符串并连接
	idString := strconv.Itoa(userID) + strconv.Itoa(questionID)

	// 使用MD5哈希函数生成哈希值
	hasher := md5.New()
	hasher.Write([]byte(idString))
	hashBytes := hasher.Sum(nil)

	// 将哈希值转换为整数
	submitID, err := strconv.ParseInt(fmt.Sprintf("%x", hashBytes), 16, 64)
	if err != nil {
		return -1, err
	}

	return submitID, nil
}

func HashPassword(from string) string {
	hash := sha256.Sum256([]byte(from))
	return hex.EncodeToString(hash[:])
}
