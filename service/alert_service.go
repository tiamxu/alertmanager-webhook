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
	"github.com/tiamxu/alertmanager-webhook/utils"
	"github.com/tiamxu/kit/log"
)

type AlertService struct{}

func NewAlertService() *AlertService {
	return &AlertService{}
}

// ProcessAlert 处理告警信息
func (s *AlertService) ProcessAlert(notification *model.AlertMessage, webhookType, templateName, webhookURL, atSomeOne, split, bot string) ([]map[string]interface{}, error) {

	// 1. 参数验证
	if err := s.validateParams(webhookType, webhookURL); err != nil {
		return nil, err
	}

	// 2. 获取模板
	template, err := s.getTemplate(webhookType, templateName)
	if err != nil {
		return nil, fmt.Errorf("template loading failed: %v", err)
	}
	notification.SetTemplate(template)

	// 3. 创建发送器
	sender, err := s.createSender(webhookType, webhookURL, bot)
	if err != nil {
		return nil, err
	}

	// 4. 处理告警
	if split == "true" {
		return s.processSplitAlerts(notification, sender, atSomeOne)
	}
	return s.processGroupedAlert(notification, sender, atSomeOne)
}

// validateParams 验证参数
func (s *AlertService) validateParams(webhookType, webhookURL string) error {
	if webhookType != "fs" && webhookType != "dd" {
		return fmt.Errorf("invalid webhook type: %s", webhookType)
	}

	if _, err := url.Parse(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %v", err)
	}

	return nil
}

// getTemplate 获取模板
func (s *AlertService) getTemplate(webhookType, templateName string) (*model.Template, error) {
	if templateName == "" {
		defaultTemplate := utils.DefaultFeishuTemplate
		if webhookType == "dd" {
			defaultTemplate = utils.DefaultDingtalkTemplate
		}
		return model.NewTemplate(defaultTemplate)
	}

	templateFile := filepath.Join("templates", templateName+".tmpl")
	return model.NewTemplate(templateFile)
}

// createSender 创建发送器
func (s *AlertService) createSender(webhookType, webhookURL, bot string) (model.MessageSender, error) {
	switch webhookType {
	case "fs":
		return &feishu.FeiShuSender{
			WebhookURL: webhookURL,
		}, nil
	case "dd":
		if !dingtalk.ValidateBot(bot) {
			return nil, fmt.Errorf("invalid bot: %s", bot)
		}
		return &dingtalk.DingTalkSender{
			WebhookURL: webhookURL,
			Bot:        bot,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported webhook type: %s", webhookType)
	}
}

// processSplitAlerts 处理分割的告警
func (s *AlertService) processSplitAlerts(notification *model.AlertMessage, sender model.MessageSender, atSomeOne string) ([]map[string]interface{}, error) {
	var messageData []map[string]interface{}

	for _, alert := range notification.Alerts {
		singleAlert := *notification
		singleAlert.Alerts = []model.Alert{alert}
		singleAlert.Status = alert.Status

		color, status := s.getAlertColorAndStatus(singleAlert)
		at := s.getAtSomeOne(atSomeOne, alert.Annotations)

		if err := s.sendAlert(&singleAlert, sender, color, status, at); err != nil {
			return nil, err
		}

		messageData = append(messageData, map[string]interface{}{
			"alert":     alert,
			"status":    status,
			"color":     color,
			"atSomeone": at,
		})
	}

	return messageData, nil
}

// processGroupedAlert 处理分组的告警
func (s *AlertService) processGroupedAlert(notification *model.AlertMessage, sender model.MessageSender, atSomeOne string) ([]map[string]interface{}, error) {
	color, status := s.getAlertColorAndStatus(*notification)
	at := s.getAtSomeOne(atSomeOne, notification.Alerts[0].Annotations)

	if err := s.sendAlert(notification, sender, color, status, at); err != nil {
		return nil, err
	}

	return []map[string]interface{}{
		{
			"alerts":    notification.Alerts,
			"status":    status,
			"color":     color,
			"atSomeone": at,
		},
	}, nil
}

// sendAlert 发送告警
func (s *AlertService) sendAlert(notification *model.AlertMessage, sender model.MessageSender, color, status, atSomeOne string) error {
	level := notification.ConvertLevelToInt()

	messageContent, err := notification.Template.Execute(notification)
	if err != nil {
		return fmt.Errorf("template execution failed: %v", err)
	}

	return sender.Send(&model.CommonMessage{
		Platform:  sender.GetPlatform(),
		Title:     notification.GroupLabels["alertname"],
		Text:      messageContent,
		Level:     level,
		Color:     color,
		Status:    status,
		AtSomeOne: atSomeOne,
	})
}

// getAtSomeOne 获取@人员
func (s *AlertService) getAtSomeOne(defaultAt string, annotations map[string]string) string {
	if at, ok := annotations["at"]; ok {
		return at
	}
	return defaultAt
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
