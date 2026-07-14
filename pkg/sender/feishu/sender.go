package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/kit/log"
)

// 全局 HTTP Client，复用连接
var httpClient = &http.Client{Timeout: 10 * time.Second}

// FeiShuSender 飞书消息发送器
type FeiShuSender struct {
	WebhookURL string
}

// NewFeishuSender 创建飞书发送器
func NewFeishuSender(webhookURL string) *FeiShuSender {
	return &FeiShuSender{WebhookURL: webhookURL}
}

// 根据告警名字@人
var alertNameUsersList = map[string]string{
	"test告警名称": "ou_1199d79525e146bad9d0a5a46a86a10f,ou_1199d79525e146bad9d0a5a46a86a10f",
}

// Send 发送卡片类型消息（主入口）
func (f *FeiShuSender) Send(message *model.CommonMessage) error {
	if message.Platform != "feishu" {
		return fmt.Errorf("invalid platform for FeiShuSender")
	}
	userOpenId := message.AtSomeOne
	SendContent := message.Text

	if userOpenId != "" {
		OpenIdtext := "<at ids=" + userOpenId + ">" + "</at>"
		SendContent += OpenIdtext
	} else {
		if _, ok := alertNameUsersList[message.Title]; ok {
			atSomeOne := alertNameUsersList[message.Title]
			OpenIdtext := "<at ids=" + atSomeOne + ">" + "</at>"
			SendContent += OpenIdtext
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
	currentTime := time.Now().Format("01-02 15:04:05")
	note := "<font color=carmine>**" + currentTime + "**</font>"

	elements := model.InteractiveMessageCardElements{
		{
			Tag:     "markdown",
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
	return f.sendRequest(msg)
}

// SendToText 发送文本消息
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
		}
		msg = model.NewPostMessage(message.Title, content)
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(f.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.WithFields(log.Fields{
			"module": "feishu",
			"action": "send_text",
			"error":  err.Error(),
		}).Error("发送文本消息失败")
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"module": "feishu",
			"action": "send_text",
			"error":  err.Error(),
		}).Error("读取响应失败")
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.WithFields(log.Fields{
			"module":   "feishu",
			"action":   "send_text",
			"status":   resp.StatusCode,
			"response": string(body),
		}).Error("文本消息发送失败")
		return fmt.Errorf("feishu webhook returned non-2xx status: %d, response: %s", resp.StatusCode, string(body))
	}
	log.WithFields(log.Fields{
		"module":   "feishu",
		"action":   "send_text",
		"response": string(body),
	}).Info("文本消息发送成功")
	return nil
}

// sendRequest 发送HTTP请求
func (f *FeiShuSender) sendRequest(msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	resp, err := httpClient.Post(f.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.WithFields(log.Fields{
			"module": "feishu",
			"action": "send",
			"error":  err.Error(),
		}).Error("发送请求失败")
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"module": "feishu",
			"action": "send",
			"error":  err.Error(),
		}).Error("读取响应失败")
		return fmt.Errorf("读取响应失败: %v", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.WithFields(log.Fields{
			"module":   "feishu",
			"action":   "send",
			"status":   resp.StatusCode,
			"response": string(body),
		}).Error("消息发送失败")
		return fmt.Errorf("feishu webhook returned non-2xx status: %d, response: %s", resp.StatusCode, string(body))
	}

	log.WithFields(log.Fields{
		"module":   "feishu",
		"action":   "send",
		"response": string(body),
	}).Info("消息发送成功")
	return nil
}

// GetPlatform 获取平台标识
func (f *FeiShuSender) GetPlatform() string {
	return "feishu"
}
