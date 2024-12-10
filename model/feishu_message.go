package model

type MsgType string

const (
	MsgTypeText        MsgType = "text"
	MsgTypePost        MsgType = "post"
	MsgTypeImage       MsgType = "image"
	MsgTypeShareChat   MsgType = "share_chat"
	MsgTypeInteractive MsgType = "interactive"
	MsgMarkDown        MsgType = "markdown"
)

// 简单文本消息
type TextMessage struct {
	MsgType MsgType            `json:"msg_type"`
	Content TextMessageContent `json:"content"`
}

type TextMessageContent struct {
	Text string `json:"text"`
}

func NewTextMessage(text string) *TextMessage {
	return &TextMessage{
		MsgType: MsgTypeText,
		Content: TextMessageContent{
			Text: text,
		},
	}
}

// 富文本消息
type PostMessage struct {
	MsgType MsgType            `json:"msg_type"`
	Content PostMessageContent `json:"content"`
}

type PostMessageContent struct {
	Post PostMessageContentPost `json:"post"`
}

type PostMessageContentPost struct {
	ZhCn PostMessageContentPostZhCn `json:"zh-CN"`
}

type PostMessageContentPostZhCn struct {
	Title   string                                `json:"title"`
	Content [][]PostMessageContentPostZhCnContent `json:"content"`
}

type PostMessageContentPostZhCnContent struct {
	Tag       string `json:"tag"`
	Text      string `json:"text,omitempty"`
	Href      string `json:"href,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	ImageKey  string `json:"image_key,omitempty"`
	FileKey   string `json:"file_key,omitempty"`
	EmojiType string `json:"emoji_type,omitempty"`
}

func NewPostMessageContentPostZhCnContent(tag, text, href, userId, userName, imageKey, fileKey, emojiType string) *PostMessageContentPostZhCnContent {
	return &PostMessageContentPostZhCnContent{
		Tag:       tag,
		Text:      text,
		Href:      href,
		UserId:    userId,
		UserName:  userName,
		ImageKey:  imageKey,
		FileKey:   fileKey,
		EmojiType: emojiType,
	}
}

func NewPostMessage(title string, content [][]PostMessageContentPostZhCnContent) *PostMessage {
	return &PostMessage{
		MsgType: MsgTypePost,
		Content: PostMessageContent{
			Post: PostMessageContentPost{
				ZhCn: PostMessageContentPostZhCn{
					Title:   title,
					Content: content,
				},
			},
		},
	}
}

// InteractiveMessage 消息卡片
type InteractiveMessage struct {
	MsgType MsgType                `json:"msg_type"`
	Card    InteractiveMessageCard `json:"card"`
}

type InteractiveMessageCard struct {
	Elements InteractiveMessageCardElements `json:"elements"`
	Header   InteractiveMessageCardHeader   `json:"header"`
}

type InteractiveMessageCardElements []struct {
	Tag     string                                `json:"tag"`
	Text    InteractiveMessageCardElementsText    `json:"text,omitempty"`
	Actions InteractiveMessageCardElementsActions `json:"actions,omitempty"`
}

type InteractiveMessageCardElementsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type InteractiveMessageCardElementsActions []struct {
	Tag   string                                    `json:"tag"`
	Text  InteractiveMessageCardElementsActionsText `json:"text"`
	Url   string                                    `json:"url"`
	Type  string                                    `json:"type"`
	Value struct {
	} `json:"value"`
}

type InteractiveMessageCardElementsActionsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type InteractiveMessageCardHeader struct {
	Title InteractiveMessageCardHeaderTitle `json:"title"`
}
type InteractiveMessageCardHeaderTitle struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

func NewInteractiveMessage(elements InteractiveMessageCardElements, header InteractiveMessageCardHeader) *InteractiveMessage {
	return &InteractiveMessage{
		MsgType: MsgTypeInteractive,
		Card: InteractiveMessageCard{
			Elements: elements,
			Header:   header,
		},
	}
}
