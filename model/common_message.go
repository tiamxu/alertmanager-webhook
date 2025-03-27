package model

import "time"

type CommonMessage struct {
	// 基础信息
	Platform string `json:"platform"` // 消息平台标识符: "feishu" 或 "dingtalk"
	Title    string `json:"title"`    // 告警标题
	Text     string `json:"text"`     // 告警内容

	// 告警相关
	Level  string `json:"level"`  // 告警等级: 信息、警告、一般严重、严重、灾难
	Status string `json:"status"` // 告警状态: 故障、恢复
	Color  string `json:"color"`  // 告警颜色: red、orange、green 等

	// 通知相关
	AtSomeOne string `json:"atSomeOne"` // @提醒的人员ID列表，逗号分隔

	// 配置相关
	Split string `json:"split"` // 是否分割告警消息: "true" 或 "false"
	Bot   string `json:"bot"`   // 机器人标识: "bot1", "bot2" 等

	// 消息类型
	Type string `json:"type"` // 消息类型: "markdown", "text" 等

	// 告警处理相关
	AlertID     string    `json:"alert_id"`     // 告警唯一标识
	HandledBy   string    `json:"handled_by"`   // 处理人
	HandledTime time.Time `json:"handled_time"` // 处理时间
	IsEscalated bool      `json:"is_escalated"` // 是否已升级
}
