package model

// OpenClawMessage OpenClaw 消息结构
type OpenClawMessage struct {
	SessionKey string                 `json:"sessionKey"`
	DingTalk   *DingTalkPayload       `json:"dingtalk"`
}

// DingTalkPayload 钉钉消息负载
type DingTalkPayload struct {
	MsgType          string      `json:"msgtype"`
	Content          string      `json:"content"`
	Sender           string      `json:"sender"`
	ConversationType string      `json:"conversationType"`
	ChatId           string      `json:"chatId"`
	MsgId            string      `json:"msgId"`
	CreateTime       int64       `json:"createTime"`
	Raw              interface{} `json:"raw"`
}

// OpenClawResponse OpenClaw 响应
type OpenClawResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ReplyRequest 回复请求
type ReplyRequest struct {
	SessionKey string `json:"sessionKey" v:"required"`
	Message    string `json:"message" v:"required"`
	Type       string `json:"type" d:"text"`
}