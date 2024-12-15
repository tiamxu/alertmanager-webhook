package model

type CommonMessage struct {
	Platform  string `json:"platform"` // 消息平台标识符，例如 "feishu" 或 "dingtalk"
	Title     string `json:"title"`
	Text      string `json:"text"`
	Type      string `json:"type"`  //类型
	Level     string `json:"level"` //告警等级
	Color     string
	Status    string
	AtSomeOne string
	Split     string
}
