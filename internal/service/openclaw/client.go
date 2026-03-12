package openclaw

import (
	"context"
	"fmt"
	"time"

	"dingtalk-bridge-go/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/glog"
)

// Client OpenClaw API 客户端
type Client struct {
	GatewayUrl string
	HookToken  string
	HookPath   string
	SessionKey string
	Timeout    int
	httpClient *gclient.Client
}

// NewClient 创建 OpenClaw 客户端
func NewClient() *Client {
	ctx := context.Background()
	
	client := &Client{
		GatewayUrl: g.Cfg().MustGet(ctx, "openclaw.gatewayUrl").String(),
		HookToken:  g.Cfg().MustGet(ctx, "openclaw.hookToken").String(),
		HookPath:   g.Cfg().MustGet(ctx, "openclaw.hookPath").String(),
		SessionKey: g.Cfg().MustGet(ctx, "openclaw.sessionKey").String(),
		Timeout:    g.Cfg().MustGet(ctx, "openclaw.timeout").Int(),
	}
	
	// 创建 HTTP 客户端
	client.httpClient = g.Client()
	client.httpClient.SetTimeout(time.Duration(client.Timeout) * time.Second)
	
	return client
}

// SendMessage 发送消息到 OpenClaw
func (c *Client) SendMessage(ctx context.Context, msg *model.DingTalkMessage, sessionKey string) (*model.OpenClawResponse, error) {
	url := fmt.Sprintf("%s%s", c.GatewayUrl, c.HookPath)
	
	if sessionKey == "" {
		sessionKey = c.SessionKey
	}
	
	// 构建请求体
	payload := &model.OpenClawMessage{
		SessionKey: sessionKey,
		DingTalk: &model.DingTalkPayload{
			MsgType:          msg.MsgType,
			Content:          msg.ExtractContent(),
			Sender:           msg.GetSenderId(),
			ConversationType: msg.ConversationType,
			ChatId:           msg.ChatId,
			MsgId:            msg.MsgId,
			CreateTime:       msg.CreateTime,
			Raw:              msg,
		},
	}
	
	glog.Info(ctx, "Sending message to OpenClaw",
		"url", url,
		"sessionKey", sessionKey,
		"sender", msg.GetSenderId(),
	)
	
	// 发送请求
	var response model.OpenClawResponse
	
	err := c.httpClient.Header(map[string]string{
		"Authorization": "Bearer " + c.HookToken,
		"Content-Type":  "application/json",
	}).PostVar(ctx, url, payload, &response)
	if err != nil {
		glog.Error(ctx, "Failed to send message to OpenClaw", err)
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	
	glog.Info(ctx, "Message sent to OpenClaw successfully")
	return &response, nil
}

// HealthCheck 检查 OpenClaw Gateway 状态
func (c *Client) HealthCheck(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/health", c.GatewayUrl)
	
	response, err := c.httpClient.Get(ctx, url)
	if err != nil {
		return false, err
	}
	defer response.Close()
	
	return response.StatusCode == 200, nil
}

// Wake 唤醒 OpenClaw
func (c *Client) Wake(ctx context.Context, text string) error {
	url := fmt.Sprintf("%s/api/wake", c.GatewayUrl)
	
	payload := g.Map{
		"text":       text,
		"sessionKey": c.SessionKey,
	}
	
	headers := map[string]string{
		"Authorization": "Bearer " + c.HookToken,
	}
	
	_, err := c.httpClient.Header(headers).Post(ctx, url, payload)
	return err
}