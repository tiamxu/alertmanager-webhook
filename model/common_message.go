package model

type CommonMessage struct {
	Platform string `json:"platform"` // 消息平台标识符，例如 "feishu" 或 "dingtalk"
	Title    string `json:"title"`
	Text     string `json:"text"`
}
