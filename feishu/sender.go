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

type FeiShuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type FeiShuSender struct {
	Name       string
	WebhookURL string
}

func (f *FeiShuSender) Send(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}
	var msg interface{}
	if message.Title == "" {
		msg = model.NewTextMessage(message.Text)
	} else {
		content := [][]model.PostMessageContentPostZhCnContent{
			{
				*model.NewPostMessageContentPostZhCnContent("markdown", message.Text, "", "", "", "", "", ""),
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
