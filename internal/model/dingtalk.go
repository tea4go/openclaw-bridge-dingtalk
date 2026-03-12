package model

// DingTalkMessage 钉钉消息结构
type DingTalkMessage struct {
	MsgType          string            `json:"msgtype" v:"required"`
	Text             *TextContent      `json:"text,omitempty"`
	Markdown         *MarkdownContent  `json:"markdown,omitempty"`
	SenderStaffId    string            `json:"senderStaffId"`
	OpenId           string            `json:"openId"`
	ConversationType string            `json:"conversationType"` // 1=单聊, 2=群聊
	ChatId           string            `json:"chatId"`
	MsgId            string            `json:"msgId"`
	CreateTime       int64             `json:"createTime"`
	AtUsers          []AtUser          `json:"atUsers,omitempty"`
	RobotCode        string            `json:"robotCode"`
}

// TextContent 文本消息内容
type TextContent struct {
	Content string `json:"content"`
}

// MarkdownContent Markdown 消息内容
type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// AtUser @用户信息
type AtUser struct {
	DingTalkId string `json:"dingtalkId"`
	StaffId    string `json:"staffId"`
}

// DingTalkWebhookRequest Webhook 请求
type DingTalkWebhookRequest struct {
	ConversationType string `json:"conversationType"`
	ChatId           string `json:"chatId"`
	MsgId            string `json:"msgId"`
	MsgType          string `json:"msgtype"`
	Text             struct {
		Content string `json:"content"`
	} `json:"text"`
	SenderStaffId string `json:"senderStaffId"`
	OpenId        string `json:"openId"`
	CreateTime    int64  `json:"createTime"`
}

// DingTalkWebhookResponse Webhook 响应
type DingTalkWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ExtractContent 提取消息内容
func (m *DingTalkMessage) ExtractContent() string {
	switch m.MsgType {
	case "text":
		if m.Text != nil {
			return m.Text.Content
		}
	case "markdown":
		if m.Markdown != nil {
			return m.Markdown.Text
		}
	case "image":
		return "[图片消息]"
	case "voice":
		return "[语音消息]"
	case "file":
		return "[文件消息]"
	}
	return ""
}

// GetSenderId 获取发送者 ID
func (m *DingTalkMessage) GetSenderId() string {
	if m.SenderStaffId != "" {
		return m.SenderStaffId
	}
	return m.OpenId
}

// IsGroupChat 是否为群聊
func (m *DingTalkMessage) IsGroupChat() bool {
	return m.ConversationType == "2"
}