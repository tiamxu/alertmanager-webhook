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
	MsgType  string `json:"msgtype"`
	Markdown Md     `json:"markdown"`
	At       At     `json:"at"`
}

func NewDDMessage(title, text string, at At) *DDMessage {
	return &DDMessage{
		MsgType: "markdown",
		Markdown: Md{
			Title: title,
			Text:  text,
		},
		At: at,
	}
}

type Text struct {
	Content string `json:"content"`
}
type DDTextMessage struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

func NewDDTextMessage(title, text string) *DDTextMessage {
	return &DDTextMessage{
		MsgType: "",
		Text: Text{
			Content: text,
		},
		At: At{
			IsAtAll: false,
		},
	}
}
