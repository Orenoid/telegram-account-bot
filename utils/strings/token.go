package strings

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
)

func GenerateToken() (string, error) {
	// 定义生成的 token 长度
	const tokenLength = 32

	// 创建一个字节数组来存储随机字节
	b := make([]byte, tokenLength)

	// 使用 crypto/rand 生成加密安全的随机字节
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	// 将字节数组编码为 base64 字符串
	token := base64.URLEncoding.EncodeToString(b)

	return token, nil
}
