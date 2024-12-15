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

// InteractiveMessage v2版本消息卡片
type InteractiveMessage struct {
	MsgType MsgType                `json:"msg_type"`
	Card    InteractiveMessageCard `json:"card"`
}

type InteractiveMessageCard struct {
	Schema   string                          `json:"schema,omitempty"`
	Elements []InteractiveMessageCardElement `json:"elements"`
	Header   InteractiveMessageCardHeader    `json:"header,omitempty"`
}

type InteractiveMessageCardElement struct {
	Tag      string                             `json:"tag"`
	Text     InteractiveMessageCardElementsText `json:"text"`
	Content  string                             `json:"content"`
	Elements []InteractiveMessageCardElement    `json:"elements"`
	// Actions InteractiveMessageCardElementsActions `json:"actions,omitempty"`
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
	Value struct{}                                  `json:"value"`
}

type InteractiveMessageCardElementsActionsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type InteractiveMessageCardHeader struct {
	Title       InteractiveMessageCardHeaderTitle `json:"title"`
	Subtitle    InteractiveMessageTagContent      `json:"subtitle"`
	TextTagList []InteractiveMessageTextTagList   `json:"text_tag_list"`
	Template    string                            `json:"template"`
}
type InteractiveMessageCardHeaderTitle struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}
type InteractiveMessageTagContent struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}
type InteractiveMessageTagText struct {
	Tag  string                       `json:"tag"`
	Text InteractiveMessageTagContent `json:"text"`
}
type InteractiveMessageTextTagList struct {
	Tag   string                       `json:"tag"`
	Text  InteractiveMessageTagContent `json:"text"`
	Color string                       `json:"color"`
}

type CardConfig struct {
	Style CardConfigStyle `json:"style"`
}
type CardConfigStyle struct {
	TextSize map[string]CardConfigTextSize `json:"text_size"`
}
type CardConfigTextSize struct {
	Default string `json:"default,omitempty"` // 使用 omitempty 忽略空值
	PC      string `json:"pc"`
	Mobile  string `json:"mobile"`
}

func NewInteractiveMessage(elements []InteractiveMessageCardElement, header InteractiveMessageCardHeader) *InteractiveMessage {
	return &InteractiveMessage{
		MsgType: MsgTypeInteractive,
		Card: InteractiveMessageCard{
			Schema:   "1.0",
			Elements: elements,
			// Body: InteractiveMessageCardBody{
			// 	Elements: elements,
			// },
			Header: header,
		},
	}
}

// 卡片2.0
type InteractiveMessageV2 struct {
	MsgType MsgType                  `json:"msg_type"`
	Card    InteractiveMessageCardV2 `json:"card"`
}

type InteractiveMessageCardV2 struct {
	Schema string                       `json:"schema"`
	Config CardConfig                   `json:"config"`
	Header InteractiveMessageCardHeader `json:"header,omitempty"`

	Body InteractiveMessageCardBody `json:"body"`
}

type InteractiveMessageCardBody struct {
	Elements InteractiveMessageCardElements `json:"elements"`
}
type InteractiveMessageCardElements []struct {
	Tag      string `json:"tag"`
	Content  string `json:"content,omitempty"`
	TextSize string `json:"text_size,omitempty"`
}

func NewInteractiveMessageV2(sytle CardConfigStyle, elements InteractiveMessageCardElements, header InteractiveMessageCardHeader) *InteractiveMessageV2 {
	return &InteractiveMessageV2{
		MsgType: MsgTypeInteractive,

		Card: InteractiveMessageCardV2{
			Schema: "2.0",
			Config: CardConfig{
				Style: sytle,
			},
			Body: InteractiveMessageCardBody{
				Elements: elements,
			},
			Header: header,
		},
	}
}
