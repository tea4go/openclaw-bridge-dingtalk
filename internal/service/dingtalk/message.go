package dingtalk

import (
	"context"
	"fmt"

	"dingtalk-bridge-go/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// Service 钉钉服务
type Service struct {
	AppSecret            string
	EnableSignatureCheck bool
}

// NewService 创建钉钉服务
func NewService() *Service {
	return &Service{
		AppSecret:            g.Cfg().MustGet(context.Background(), "dingtalk.appSecret").String(),
		EnableSignatureCheck: g.Cfg().MustGet(context.Background(), "dingtalk.enableSignatureCheck").Bool(),
	}
}

// ProcessMessage 处理钉钉消息
func (s *Service) ProcessMessage(ctx context.Context, msg *model.DingTalkMessage) error {
	glog.Info(ctx, "Processing DingTalk message",
		"sender", msg.GetSenderId(),
		"type", msg.MsgType,
		"content", msg.ExtractContent(),
	)

	// 这里可以添加消息过滤、转换等逻辑
	
	return nil
}

// VerifySignature 验证消息签名
func (s *Service) VerifySignature(ctx context.Context, sign, timestamp string) bool {
	if !s.EnableSignatureCheck {
		glog.Debug(ctx, "Signature check disabled")
		return true
	}

	if sign == "" || timestamp == "" {
		glog.Warning(ctx, "Missing signature or timestamp")
		return false
	}

	// 检查时间戳
	if !CheckTimestamp(timestamp) {
		glog.Warning(ctx, "Invalid timestamp", timestamp)
		return false
	}

	// 验证签名
	isValid := VerifyOutgoingSignature(sign, timestamp, s.AppSecret)
	
	if !isValid {
		glog.Warning(ctx, "Invalid signature")
	} else {
		glog.Debug(ctx, "Signature verified")
	}
	
	return isValid
}

// BuildTextMessage 构建文本消息
func BuildTextMessage(content string) map[string]interface{} {
	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
}

// BuildMarkdownMessage 构建 Markdown 消息
func BuildMarkdownMessage(title, text string) map[string]interface{} {
	return map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}
}

// GetSessionKey 生成会话 Key
func GetSessionKey(senderId, chatId string) string {
	if chatId != "" {
		return fmt.Sprintf("dingtalk:group:%s", chatId)
	}
	return fmt.Sprintf("dingtalk:dm:%s", senderId)
}