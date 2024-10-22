package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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
