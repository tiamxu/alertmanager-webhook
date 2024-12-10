package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
)

type FeiShuSender struct {
	Name       string
	WebhookURL string
}

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

// 发送卡片类型消息
func (f *FeiShuSender) Send(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}
	var msg interface{}
	elements := model.InteractiveMessageCardElements{
		{
			Tag:     "markdown",
			Content: message.Text,
			// Content: "# 一级标题\n## 二级标题\n### 三级标题\n#### 四级标题\n##### 五级标题\n###### 六级标题",
		},
		{
			Tag: "hr",
		},
	}

	headers := model.InteractiveMessageCardHeader{
		Title: model.InteractiveMessageCardHeaderTitle{
			Content: "测试告警",
			Tag:     "plain_text",
		},
	}
	msg = model.NewInteractiveMessage(elements, headers)
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Printf("payload:%s\n", string(payload))
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
