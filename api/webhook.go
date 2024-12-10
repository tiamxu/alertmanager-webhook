package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/config"
	"github.com/tiamxu/alertmanager-webhook/dingtalk"
	"github.com/tiamxu/alertmanager-webhook/feishu"
	"github.com/tiamxu/alertmanager-webhook/interfaces"
	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
)

type FeiShuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func HandlerWebhook(c *gin.Context) {
	var notification model.AlertMessage
	err := c.BindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	webhookType := c.Query("type")
	templateName := c.Query("tpl")
	fsURL := c.Query("fsurl")
	if webhookType != "fs" || templateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing parameters"})
		return
	}
	// templateName := notification.GetTemplateName()
	templateFile := filepath.Join("templates", templateName+".tmpl")
	alertTemplate, err := model.NewTemplate(templateFile)
	if err != nil {
		log.Errorf("Failed to load template file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "template loading failed"})
		return
	}
	notification.SetTemplate(alertTemplate)
	messageContent, err := notification.Template.Execute(notification)
	// messageContext, err := alertTemplate.Execute(notification)
	if err != nil {
		log.Errorf("Failed to execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "template execution failed"})
		return
	}
	// 构造 JSON 请求体
	feishuMessage := map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": "Prometheus Alert",
					"content": [][]map[string]interface{}{
						{
							{
								"tag":  "text",
								"text": messageContent,
							},
						},
					},
				},
			},
		},
	}
	jsonPayload, err := json.Marshal(feishuMessage)
	if err != nil {
		log.Errorf("Failed to marshal JSON payload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal JSON"})
		return
	}
	parsedURL, err := url.Parse(fsURL)
	if err != nil {
		log.Errorf("Failed to parse FeiShu webhook URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fsurl parameter"})
		return
	}
	// for i := range notification.Alerts {
	// 	notification.Alerts[i].Annotations["text"] = messageContent
	// }
	// commonMsg := &model.CommonMessage{
	// 	Platform: config.AppConfig.AlertType,
	// 	Title:    notification.GroupLabels["alertname"],
	// 	Text:     messageContent,
	// }
	commonMsg := &model.CommonMessage{
		Text: string(jsonPayload),
		// Title:    notification.GroupLabels["alertname"],
		Platform: config.AppConfig.AlertType,
	}
	// sender := getSender(commonMsg)
	sender := &feishu.FeiShuSender{
		Name:       "robot1",
		WebhookURL: parsedURL.String(),
	}
	//fmt.Println(msg)
	// if sender == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported platform"})
	// 	log.Infof("getSender:%s", err)

	// 	return
	// }

	err = sender.Send(commonMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successful send alert notification!"})
}

func getSender(notification *model.CommonMessage) interfaces.MessageSender {
	log.Infof("Attempting to get sender for platform: %s", notification.Platform)
	switch notification.Platform {
	case "dingtalk":
		if config.AppConfig.OpenDingding == 1 {
			return &dingtalk.DingTalkSender{WebhookURL: config.AppConfig.Dingtalk.WebhookURL}
		}
	case "feishu":
		if config.AppConfig.OpenFeishu == 1 {
			return &feishu.FeiShuSender{WebhookURL: config.AppConfig.Feishu.WebhookURL}

		}
	default:
		log.Warnf("Unsupported platform specified: %s", notification.Platform)
	}

	return nil
}

// func buildMessageText(notification model.AlertMessage) string {
// 	var buffer strings.Builder
// 	buffer.WriteString(fmt.Sprintf("通知组%s,状态[%s]\n告警项\n\n", notification.GroupKey, notification.Status))

// 	for _, alert := range notification.Alerts {
// 		buffer.WriteString(fmt.Sprintf("摘要：%s\n详情: %s\n", alert.Annotations["summary"], alert.Annotations["description"]))
// 		buffer.WriteString(fmt.Sprintf("开始时间: %s\n\n", alert.StartsAt.Format("2023-12-01 15:04:05")))
// 	}

// 	return buffer.String()
// }
