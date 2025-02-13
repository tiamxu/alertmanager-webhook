package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/tiamxu/alertmanager-webhook/dingtalk"
	"github.com/tiamxu/alertmanager-webhook/feishu"
	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/kit/log"
)

type AlertService struct{}

func NewAlertService() *AlertService {
	return &AlertService{}
}

func (s *AlertService) ProcessAlert(notification *model.AlertMessage, webhookType, templateName, webhookURL, atSomeOne, split string) ([]map[string]interface{}, error) {
	if webhookType != "fs" && webhookType != "dd" || templateName == "" {
		return nil, fmt.Errorf("invalid or missing parameters")
	}

	// 转告警级别为中文并设置消息颜色和状态
	level := notification.ConvertLevelToInt()
	color, status := s.getAlertColorAndStatus(*notification)

	// 加载模板文件
	templateFile := filepath.Join("templates", templateName+".tmpl")
	alertTemplate, err := model.NewTemplate(templateFile)
	if err != nil {
		return nil, fmt.Errorf("template loading failed: %v", err)
	}
	notification.SetTemplate(alertTemplate)

	// 解析  webhook URL
	_, err = url.Parse(webhookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook url parameter: %v", err)
	}

	var messageData []map[string]interface{}

	// 根据不同类型创建不同的发送器
	var sender model.MessageSender
	var platform string

	switch webhookType {
	case "fs":
		platform = "feishu"
		sender = &feishu.FeiShuSender{
			WebhookURL: webhookURL,
		}
	case "dd":
		platform = "dingtalk"
		// 从 URL 中解析 secret
		secret := ""
		if u, err := url.Parse(webhookURL); err == nil {
			secret = u.Query().Get("secret")
		}
		sender = &dingtalk.DingTalkSender{
			WebhookURL: webhookURL,
			Secret:     secret,
		}
	}

	// 构建和发送消息
	if split == "true" {
		for _, alert := range notification.Alerts {
			notification.Alerts = []model.Alert{alert}
			notification.Status = alert.Status
			color, status := s.getAlertColorAndStatus(*notification)
			at := atSomeOne
			if atInner, ok := alert.Annotations["at"]; ok {
				at = atInner
			}
			_, err := s.sendAlertMessage(
				level,
				color,
				status,
				at,
				notification.GroupLabels["alertname"],
				notification,
				alertTemplate,
				sender,
				platform,
			)
			if err != nil {
				return nil, fmt.Errorf("message sending failed: %v", err)
			}
			messageData = append(messageData, map[string]interface{}{
				"alert":     alert,
				"status":    status,
				"color":     color,
				"atSomeone": at,
			})
		}
	} else {
		at := atSomeOne
		if atInner, ok := notification.Alerts[0].Annotations["at"]; ok {
			at = atInner
		}
		_, err := s.sendAlertMessage(
			level,
			color,
			status,
			at,
			notification.GroupLabels["alertname"],
			notification,
			alertTemplate,
			sender,
			platform,
		)
		if err != nil {
			return nil, fmt.Errorf("message sending failed: %v", err)
		}
		messageData = append(messageData, map[string]interface{}{
			"alerts":    notification.Alerts,
			"status":    status,
			"color":     color,
			"atSomeone": atSomeOne,
		})
	}

	return messageData, nil
}

func (s *AlertService) sendAlertMessage(level, color, status, atSomeOne, title string, message interface{}, tmpl *model.Template, sender model.MessageSender, platform string) (string, error) {
	messageContent, err := tmpl.Execute(message)
	if err != nil {
		return "", fmt.Errorf("template execution failed: %v", err)
	}

	commonMsg := &model.CommonMessage{
		Platform:  platform,
		Title:     title,
		Text:      messageContent,
		Level:     level,
		Color:     color,
		Status:    status,
		AtSomeOne: atSomeOne,
	}

	if err := sender.Send(commonMsg); err != nil {
		return "", fmt.Errorf("send message failed: %v", err)
	}

	return messageContent, nil
}

func (s *AlertService) getAlertColorAndStatus(notification model.AlertMessage) (string, string) {
	content, err := json.Marshal(notification)
	if err != nil {
		log.Errorf("getAlertColorAndStatus Error marshalling JSON: %v", err)
		return "red", "故障" // 默认返回红色和"故障"状态
	}
	contentStr := strings.ToLower(string(content))
	switch {
	case strings.Contains(contentStr, "resolved") && strings.Contains(contentStr, "firing"):
		return "orange", "故障"
	case strings.Contains(contentStr, "resolved"):
		return "green", "恢复"
	default:
		return "red", "故障"
	}
}
