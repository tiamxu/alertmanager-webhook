package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tiamxu/alertmanager-webhook/model"
)

type DingTalkSender struct {
	Name       string
	WebhookURL string
}

func (d *DingTalkSender) Send(message *model.CommonMessage) error {
	if message.Platform != "dingtalk" {
		return fmt.Errorf("invalid platform for DingTalkSender")
	}
	data := model.DDMessage{
		Msgtype: "markdown",
		Markdown: model.Md{
			Title: message.Title,
			Text:  message.Text,
		},
		At: model.At{
			IsAtAll: false,
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("DingTalk response: %s\n", body)
	return nil
}
