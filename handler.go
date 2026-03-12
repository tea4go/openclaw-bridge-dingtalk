package main

import (
	"net/http"

	"dingtalk-bridge-go/internal/model"
	"dingtalk-bridge-go/internal/service/dingtalk"
	"dingtalk-bridge-go/internal/service/openclaw"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

// 全局服务实例
var (
	dingtalkService *dingtalk.Service
	openclawClient  *openclaw.Client
)

func init() {
	dingtalkService = dingtalk.NewService()
	openclawClient = openclaw.NewClient()
}

// DingTalkWebhookHandler 钉钉 Webhook 处理器
func DingTalkWebhookHandler(r *ghttp.Request) {
	ctx := r.Context()
	
	// 获取签名头
	sign := r.Header.Get("sign")
	timestamp := r.Header.Get("timestamp")
	
	glog.Info(ctx, "Received DingTalk webhook",
		"sign", sign,
		"timestamp", timestamp,
	)
	
	// 验证签名
	if !dingtalkService.VerifySignature(ctx, sign, timestamp) {
		glog.Warning(ctx, "Signature verification failed")
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Invalid signature",
		})
		return
	}
	
	// 解析消息
	var msg model.DingTalkMessage
	if err := r.Parse(&msg); err != nil {
		glog.Error(ctx, "Failed to parse message", err)
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Invalid message format",
		})
		return
	}
	
	// 只处理文本消息
	if msg.MsgType != "text" {
		glog.Info(ctx, "Ignoring non-text message", "type", msg.MsgType)
		r.Response.WriteJson(g.Map{
			"success": true,
			"message": "Non-text message ignored",
		})
		return
	}
	
	content := msg.ExtractContent()
	if content == "" {
		r.Response.WriteJson(g.Map{
			"success": true,
			"message": "Empty message",
		})
		return
	}
	
	glog.Info(ctx, "Processing message",
		"sender", msg.GetSenderId(),
		"content", content,
		"isGroup", msg.IsGroupChat(),
	)
	
	// 处理消息
	if err := dingtalkService.ProcessMessage(ctx, &msg); err != nil {
		glog.Error(ctx, "Failed to process message", err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Failed to process message",
		})
		return
	}
	
	// 生成会话 Key
	sessionKey := dingtalk.GetSessionKey(msg.GetSenderId(), msg.ChatId)
	
	// 发送到 OpenClaw
	_, err := openclawClient.SendMessage(ctx, &msg, sessionKey)
	if err != nil {
		glog.Error(ctx, "Failed to send to OpenClaw", err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Failed to forward message",
		})
		return
	}
	
	// 返回成功响应
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Message received and processing",
	})
}

// ReplyHandler 回复处理器 (OpenClaw 回调)
func ReplyHandler(r *ghttp.Request) {
	ctx := r.Context()
	
	var req model.ReplyRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}
	
	glog.Info(ctx, "Received reply",
		"sessionKey", req.SessionKey,
		"type", req.Type,
		"message", req.Message,
	)
	
	// TODO: 实现回复到钉钉的逻辑
	// 需要存储 sender 和 webhook URL 的映射关系
	
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Reply received",
	})
}