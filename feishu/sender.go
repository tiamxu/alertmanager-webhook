package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
)

type At struct {
	AlertName string `json:"alertname"`
	AtSomeOne string `json:"atSomeOne"`
}

type FeiShuSender struct {
	Name       string
	WebhookURL string
}

// 根据告警名字@人
var alertNameUsersList = map[string]string{
	"cpu告警":  "ou_1199d79525e146bad9d0a5a46a86a10f,ou_1199d79525e146bad9d0a5a46a86a10f",
	"服务异常告警": "ou_1199d79525e146bad9d0a5a46a86a10f",
}

// 富文本
func (f *FeiShuSender) SendToText(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}
	var msg interface{}
	if message.Title == "" {
		msg = model.NewTextMessage(message.Text)
	} else {
		content := [][]model.PostMessageContentPostZhCnContent{
			{
				*model.NewPostMessageContentPostZhCnContent("text", message.Text, "", "", "", "", "", ""),
			},
			// {
			// 	*model.NewPostMessageContentPostZhCnContent("a", "点击查看", "http://www.baidu.com", "", "", "", "", ""),
			// },
		}
		msg = model.NewPostMessage(message.Title, content)
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(f.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Errorln("[feishu]", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("[feishu]", err.Error())
	}
	log.Infoln("[feishu]", string(body))
	return nil

}

// 发送卡片类型消息schema=1.0版本
func (f *FeiShuSender) Send(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}

	var color, userOpenId string
	if strings.Count(message.Text, "resolved") > 0 && strings.Count(message.Text, "firing") > 0 {
		color = "orange"
	} else if strings.Count(message.Text, "resolved") > 0 {
		color = "green"
	} else {
		color = "red"
	}

	SendContent := message.Text

	if userOpenId != "" {
		OpenIds := strings.Split(userOpenId, ",")
		OpenIdtext := ""
		for _, OpenId := range OpenIds {
			OpenIdtext += "<at user_id=" + OpenId + " id=" + OpenId + " email=" + OpenId + "></at>"
		}
		SendContent += OpenIdtext
	}
	// SendContent += "<at id=7a22d6ab></at>"
	// SendContent += "<at id=l.xu@unipal.com.cn></at>"
	var msg interface{}
	// currentTime := time.Now().Format("2006-01-02 15:04:05") // 使用标准布局格式化时间
	headers := model.InteractiveMessageCardHeader{
		Title: model.InteractiveMessageCardHeaderTagContent{
			Content: message.Title,
			Tag:     "plain_text",
		},
		// 标题主题颜色。支持 "blue"|"wathet"|"tuiquoise"|"green"|"yellow"|"orange"|"red"|"carmine"|"violet"|"purple"|"indigo"|"grey"|"default"。默认值 default。
		Template: color,
	}

	elements := []model.InteractiveMessageCardElement{
		// {
		// 	Tag:     "markdown",
		// 	Content: "### 测试",
		// },
		{
			Tag: "div",
			Text: model.InteractiveMessageCardElementsText{
				Tag:     "lark_md",
				Content: SendContent,
			},
		},
		{
			Tag: "hr",
		},
		// {
		// 	Tag: "note",
		// 	Elements: []model.InteractiveMessageCardElement{
		// 		{
		// 			Tag:     "lark_md",
		// 			Content: fmt.Sprintf("生成时间: %s", currentTime),
		// 		},
		// 	},
		// },
	}

	msg = model.NewInteractiveMessage(elements, headers)

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// fmt.Printf("payload:%s\n", string(payload))
	resp, err := http.Post(f.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Errorln("[feishu]", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("[feishu]", err.Error())
	}
	log.Infoln("[feishu]", string(body))
	return nil

}

// 发送卡片类型消息schema=2.0版本
func (f *FeiShuSender) SendV2(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}
	userOpenId := message.AtSomeOne
	SendContent := message.Text

	if userOpenId != "" {
		// OpenIds := strings.Split(userOpenId, ",")
		OpenIdtext := "<at ids=" + userOpenId + ">" + "</at>"
		// for _, OpenId := range OpenIds {
		// 	OpenIdtext += "<at user_id=" + OpenId + " id=" + OpenId + " email=" + OpenId + "></at>"
		// }
		SendContent += OpenIdtext
	} else {

		for alertName, atSomeOne := range alertNameUsersList {
			if alertName == message.Title {
				OpenIdtext := "<at ids=" + atSomeOne + ">" + "</at>"
				SendContent += OpenIdtext
			}
		}
	}
	var msg interface{}
	style := model.CardConfigStyle{
		TextSize: map[string]model.CardConfigTextSize{
			"cus-0": {
				Default: "medium",
				PC:      "x-large",
				Mobile:  "large",
			},
		},
	}
	headers := model.InteractiveMessageCardHeader{
		Title: model.InteractiveMessageCardHeaderTagContent{
			Content: message.Title,
			Tag:     "plain_text",
		},
		// Subtitle: model.InteractiveMessageTagContent{
		// 	Tag:     "plain_text",
		// 	Content: "恢复",
		// },
		TextTagList: []model.InteractiveMessageHeaderTextTagList{
			{
				Tag: "text_tag",
				Text: model.InteractiveMessageCardHeaderTagContent{
					Tag:     "plain_text",
					Content: message.Level,
				},
				Color: message.Color,
			},
			{
				Tag: "text_tag",
				Text: model.InteractiveMessageCardHeaderTagContent{
					Tag:     "plain_text",
					Content: message.Status,
				},
				Color: message.Color,
			},
		},
		Template: message.Color,
	}
	currentTime := time.Now().Format("01-02 15:04:05") // 使用标准布局格式化时间
	note := "<font color=carmine>**" + currentTime + "**</font>"

	elements := model.InteractiveMessageCardElements{

		{
			Tag: "markdown",
			// TextSize: "cus-0",
			Content: SendContent,
		},
		{
			Tag: "hr",
		},
		{
			Tag:      "markdown",
			TextSize: "notation",
			Content:  note,
		},
	}

	msg = model.NewInteractiveMessageV2(style, elements, headers)

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// fmt.Printf("payload:%s\n", string(payload))
	resp, err := http.Post(f.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Errorln("[feishuv2]", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("[feishuv2]", err.Error())
	}
	log.Infoln("[feishuv2]", string(body))
	return nil

}

func (f *FeiShuSender) NewFeiShuSender() {
	return
}
