package model

// 定义钉钉消息类型

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type Md struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DDMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown Md     `json:"markdown"`
	At       At     `json:"at"`
}
