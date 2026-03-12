package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// Crypto 钉钉加密工具
type Crypto struct {
	Token          string
	EncodingAESKey string
	AppKey         string
}

// NewCrypto 创建加密工具实例
func NewCrypto(token, encodingAESKey, appKey string) *Crypto {
	return &Crypto{
		Token:          token,
		EncodingAESKey: encodingAESKey,
		AppKey:         appKey,
	}
}

// VerifyOutgoingSignature 验证 Outgoing 机器人签名
// 钉钉 Outgoing 机器人使用 HMAC-SHA256 签名
func VerifyOutgoingSignature(sign, timestamp, appSecret string) bool {
	str := fmt.Sprintf("%s\n%s", timestamp, appSecret)
	
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(str))
	expectedSign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	
	return sign == expectedSign
}

// GenerateOutgoingSign 生成 Outgoing 签名
func GenerateOutgoingSign(timestamp, appSecret string) string {
	str := fmt.Sprintf("%s\n%s", timestamp, appSecret)
	
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(str))
	
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifyUrlSignature 验证回调 URL 签名
func (c *Crypto) VerifyUrlSignature(signature, timestamp, nonce, encrypt string) bool {
	str := sortAndJoin(c.Token, timestamp, nonce, encrypt)
	hash := sha256.Sum256([]byte(str))
	expectedSign := fmt.Sprintf("%x", hash)
	
	return signature == expectedSign
}

// VerifyMessageSignature 验证消息签名
func (c *Crypto) VerifyMessageSignature(signature, timestamp, nonce, encrypt string) bool {
	return c.VerifyUrlSignature(signature, timestamp, nonce, encrypt)
}

// GenerateSignature 生成签名
func (c *Crypto) GenerateSignature(timestamp, nonce, encrypt string) string {
	str := sortAndJoin(c.Token, timestamp, nonce, encrypt)
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

// sortAndJoin 排序并连接字符串
func sortAndJoin(items ...string) string {
	// 简单的冒泡排序
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i] > items[j] {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return strings.Join(items, "")
}

// CheckTimestamp 检查时间戳是否有效 (防止重放攻击)
// 允许 5 分钟的时间差
func CheckTimestamp(timestamp string) bool {
	ts, err := parseInt64(timestamp)
	if err != nil {
		return false
	}
	
	now := time.Now().UnixMilli()
	diff := now - ts
	
	// 允许 5 分钟的时间差
	return diff >= -300000 && diff <= 300000
}

func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}