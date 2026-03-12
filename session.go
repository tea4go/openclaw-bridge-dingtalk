package main

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/glog"
)

// SessionManager 会话管理器
type SessionManager struct {
	sessions map[string]*Session
	maxAge   int64 // 秒
}

// Session 会话信息
type Session struct {
	Key          string
	SenderId     string
	ChatId       string
	CreatedAt    int64
	LastActivity int64
}

// NewSessionManager 创建会话管理器
func NewSessionManager() *SessionManager {
	maxAge := g.Cfg().MustGet(context.Background(), "session.maxAge").Int64()
	if maxAge == 0 {
		maxAge = 86400 // 默认 24 小时
	}
	
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		maxAge:   maxAge,
	}
	
	// 启动清理任务
	go sm.cleanupTask()
	
	return sm
}

// GetOrCreate 获取或创建会话
func (sm *SessionManager) GetOrCreate(senderId, chatId string) *Session {
	key := sm.generateKey(senderId, chatId)
	
	if session, exists := sm.sessions[key]; exists {
		session.LastActivity = gtime.Now().Timestamp()
		return session
	}
	
	now := gtime.Now().Timestamp()
	session := &Session{
		Key:          key,
		SenderId:     senderId,
		ChatId:       chatId,
		CreatedAt:    now,
		LastActivity: now,
	}
	
	sm.sessions[key] = session
	glog.Info(context.Background(), "Created new session", "key", key)
	
	return session
}

// Get 获取会话
func (sm *SessionManager) Get(key string) (*Session, bool) {
	session, exists := sm.sessions[key]
	if exists {
		session.LastActivity = gtime.Now().Timestamp()
	}
	return session, exists
}

// generateKey 生成会话 Key
func (sm *SessionManager) generateKey(senderId, chatId string) string {
	if chatId != "" {
		return "dingtalk:group:" + chatId
	}
	return "dingtalk:dm:" + senderId
}

// cleanupTask 清理过期会话任务
func (sm *SessionManager) cleanupTask() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		sm.cleanup()
	}
}

// cleanup 清理过期会话
func (sm *SessionManager) cleanup() {
	ctx := context.Background()
	now := gtime.Now().Timestamp()
	count := 0
	
	for key, session := range sm.sessions {
		if now-session.LastActivity > sm.maxAge {
			delete(sm.sessions, key)
			count++
		}
	}
	
	if count > 0 {
		glog.Info(ctx, "Cleaned up expired sessions", "count", count)
	}
}

// Stats 获取统计信息
func (sm *SessionManager) Stats() map[string]interface{} {
	return map[string]interface{}{
		"total_sessions": len(sm.sessions),
		"max_age":        sm.maxAge,
	}
}