package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/tiamxu/alertmanager-webhook/config"
	"github.com/tiamxu/alertmanager-webhook/feishu"
	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/alertmanager-webhook/pkg/e"
)

type AlertService struct{}

func NewAlertService() *AlertService {
	return &AlertService{}
}

func (s *AlertService) ProcessAlert(notification *model.AlertMessage, webhookType, templateName, fsURL, atSomeOne, split string) (model.Response, error) {
	code := e.SUCCESS
	if webhookType != "fs" || templateName == "" {
		code := e.InvalidParams
		// return nil, fmt.Errorf("invalid or missing parameters")
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
		}, fmt.Errorf("invalid or missing parameters")
	}

	// 转告警级别为中文并设置消息颜色和状态
	level := notification.ConvertLevelToInt()
	color, status := s.getAlertColorAndStatus(*notification)

	// 加载模板文件
	templateFile := filepath.Join("templates", templateName+".tmpl")
	alertTemplate, err := model.NewTemplate(templateFile)
	if err != nil {
		code := e.ErrorTemplateLoad
		// return nil, fmt.Errorf("template loading failed: %v", err)
		return model.Response{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}, fmt.Errorf("template loading failed: %v", err)
	}
	notification.SetTemplate(alertTemplate)

	// 解析 FeiShu webhook URL
	parsedURL, err := url.Parse(fsURL)
	if err != nil {
		// return nil, fmt.Errorf("invalid fsurl parameter: %v", err)
		code := e.InvalidParams
		return model.Response{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}, fmt.Errorf("invalid fsurl parameter: %v", err)
	}

	var messageData []map[string]interface{}

	// 构建和发送消息
	if split == "true" && atSomeOne == "" {
		for _, alert := range notification.Alerts {
			notification.Alerts = []model.Alert{alert}
			notification.Status = alert.Status
			color, status := s.getAlertColorAndStatus(*notification)
			at := atSomeOne
			if atInner, ok := alert.Annotations["at"]; ok {
				at = atInner
			}
			_, err := s.sendAlertMessage(parsedURL.String(), level, color, status, at, alert.Labels["alertname"], notification, alertTemplate)
			if err != nil {
				code := e.ERROR
				return model.Response{
					Code:  code,
					Msg:   e.GetMsg(code),
					Error: err.Error(),
				}, fmt.Errorf("message sending failed: %v", err)

				// return nil, fmt.Errorf("message sending failed: %v", err)
			}
			messageData = append(messageData, map[string]interface{}{
				"alert":     alert,
				"status":    status,
				"color":     color,
				"atSomeone": at,
			})
		}
	} else {
		_, err := s.sendAlertMessage(parsedURL.String(), level, color, status, atSomeOne, notification.GroupLabels["alertname"], notification, alertTemplate)
		if err != nil {
			code := e.ERROR
			return model.Response{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: err.Error(),
			}, fmt.Errorf("message sending failed: %v", err)
			// return nil, fmt.Errorf("message sending failed: %v", err)
		}
		messageData = append(messageData, map[string]interface{}{
			"alerts":    notification.Alerts,
			"status":    status,
			"color":     color,
			"atSomeone": atSomeOne,
		})
	}

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: messageData,
	}, nil
}

func (s *AlertService) sendAlertMessage(fsurl, level, color, status, atSomeOne, title string, message interface{}, tmpl *model.Template) (string, error) {
	messageContent, err := tmpl.Execute(message)
	if err != nil {
		return "", err
	}

	commonMsg := &model.CommonMessage{
		Platform:  config.AppConfig.AlertType,
		Title:     title,
		Text:      messageContent,
		Level:     level,
		Color:     color,
		Status:    status,
		AtSomeOne: atSomeOne,
	}

	sender := &feishu.FeiShuSender{
		WebhookURL: fsurl,
	}
	if err := sender.SendV2(commonMsg); err != nil {
		return "", err
	}

	return messageContent, nil
}

func (s *AlertService) getAlertColorAndStatus(notification model.AlertMessage) (string, string) {
	content, err := json.Marshal(notification)
	if err != nil {
		log.Errorf("getAlertColorAndStatus Error marshalling JSON: %v", err)
		return "red", "故障" // 默认返回红色和“故障”状态
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
