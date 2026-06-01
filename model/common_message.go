package model

// CommonMessage 通用消息结构（跨平台）
type CommonMessage struct {
	Platform  string `json:"platform"` // 消息平台标识符: "feishu" 或 "dingtalk"
	Title     string `json:"title"`    // 消息标题
	Text      string `json:"text"`     // 消息内容
	Level     string `json:"level"`    // 告警等级: 信息、警告、一般严重、严重、灾难
	Status    string `json:"status"`   // 告警状态: 故障、恢复
	Color     string `json:"color"`   // 告警颜色: red、orange、green 等
	AtSomeOne string `json:"atSomeOne"` // @提醒的人员ID列表，逗号分隔
}
