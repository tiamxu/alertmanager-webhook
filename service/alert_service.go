package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/alertmanager-webhook/pkg/sender"
	"github.com/tiamxu/kit/log"
)

type AlertService struct{}

func NewAlertService() *AlertService {
	return &AlertService{}
}

// ProcessAlert 处理告警信息
func (s *AlertService) ProcessAlert(notification *model.AlertMessage, webhookType, templateName, webhookURL, atSomeOne, split, bot string) ([]map[string]interface{}, error) {

	// 1. 参数验证
	if notification == nil {
		return nil, fmt.Errorf("notification is nil")
	}
	if len(notification.Alerts) == 0 {
		return nil, fmt.Errorf("alerts is empty")
	}
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
	alertSender, err := s.createSender(webhookType, webhookURL, bot)
	if err != nil {
		return nil, err
	}

	// 4. 处理告警
	if split == "true" {
		return s.processSplitAlerts(notification, alertSender, atSomeOne)
	}
	return s.processGroupedAlert(notification, alertSender, atSomeOne)
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
		templateName = "feishu"
		if webhookType == "dd" {
			templateName = "dingtalk"
		}
	}

	templateFile := filepath.Join("templates", templateName+".tmpl")
	return model.NewTemplate(templateFile)
}

// createSender 创建发送器
func (s *AlertService) createSender(webhookType, webhookURL, bot string) (sender.MessageSender, error) {
	return sender.NewSender(webhookType, webhookURL, bot)
}

// processSplitAlerts 处理分割的告警（并发发送）
func (s *AlertService) processSplitAlerts(notification *model.AlertMessage, alertSender sender.MessageSender, atSomeOne string) ([]map[string]interface{}, error) {
	var (
		wg          sync.WaitGroup
		messageData []map[string]interface{}
		mu          sync.Mutex
		errChan     = make(chan error, len(notification.Alerts))
	)

	for _, alert := range notification.Alerts {
		wg.Add(1)
		go func(alert model.Alert) {
			defer wg.Done()

			singleAlert := *notification
			singleAlert.Alerts = []model.Alert{alert}
			singleAlert.Status = alert.Status

			color, status := s.getAlertColorAndStatus(singleAlert)
			at := s.getAtSomeOne(atSomeOne, alert.Annotations)

			if err := s.sendAlert(&singleAlert, alertSender, color, status, at); err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			messageData = append(messageData, map[string]interface{}{
				"alert":     alert,
				"status":    status,
				"color":     color,
				"atSomeone": at,
			})
			mu.Unlock()
		}(alert)
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return messageData, err
		}
	}

	return messageData, nil
}

// processGroupedAlert 处理分组的告警
func (s *AlertService) processGroupedAlert(notification *model.AlertMessage, alertSender sender.MessageSender, atSomeOne string) ([]map[string]interface{}, error) {
	color, status := s.getAlertColorAndStatus(*notification)
	at := s.getAtSomeOne(atSomeOne, notification.Alerts[0].Annotations)

	if err := s.sendAlert(notification, alertSender, color, status, at); err != nil {
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
func (s *AlertService) sendAlert(notification *model.AlertMessage, alertSender sender.MessageSender, color, status, atSomeOne string) error {
	level := notification.ConvertLevelToInt()

	messageContent, err := notification.Template.Execute(notification)
	if err != nil {
		log.WithFields(log.Fields{
			"module":    "alert_service",
			"action":    "send_alert",
			"alertname": notification.GroupLabels["alertname"],
			"error":     err.Error(),
		}).Error("模板渲染失败")
		return fmt.Errorf("template execution failed: %v", err)
	}

	if err := alertSender.Send(&model.CommonMessage{
		Platform:  alertSender.GetPlatform(),
		Title:     notification.GroupLabels["alertname"],
		Text:      messageContent,
		Level:     level,
		Color:     color,
		Status:    status,
		AtSomeOne: atSomeOne,
	}); err != nil {
		log.Errorf("发送告警失败 [platform=%s, alertname=%s]: %v", alertSender.GetPlatform(), notification.GroupLabels["alertname"], err)
		return err
	}

	return nil
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
		return "red", "故障"
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
