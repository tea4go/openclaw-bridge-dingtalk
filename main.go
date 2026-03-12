package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	ctx := context.Background()

	// 初始化会话管理器
	sessionManager := NewSessionManager()
	_ = sessionManager

	// 创建 HTTP 服务器
	s := g.Server()

	// 配置路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareCORS)
		
		// 健康检查
		group.GET("/health", HealthHandler)
		group.GET("/", IndexHandler)
		
		// 钉钉 Webhook
		group.POST("/webhook/dingtalk", DingTalkWebhookHandler)
		
		// 回复接口 (OpenClaw 回调)
		group.POST("/api/reply", ReplyHandler)
	})

	// 优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-c
		glog.Info(ctx, "Shutting down server...")
		s.Shutdown()
	}()

	// 启动服务器
	glog.Info(ctx, "Starting DingTalk-OpenClaw Bridge Server...")
	glog.Info(ctx, "Server listening on", s.GetListenedAddress())
	glog.Info(ctx, "Webhook endpoint: POST /webhook/dingtalk")
	
	s.Run()
}

// HealthHandler 健康检查
func HealthHandler(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"status":    "ok",
		"timestamp": gtime.Now().TimestampMilli(),
	})
}

// IndexHandler 服务状态
func IndexHandler(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"name":    "DingTalk-OpenClaw Bridge",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": g.Map{
			"webhook": "POST /webhook/dingtalk",
			"health":  "GET /health",
			"reply":   "POST /api/reply",
		},
	})
}