package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/config"
	"github.com/tiamxu/alertmanager-webhook/feishu"
	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
)

func PrometheusAlert(c *gin.Context) {
	var notification model.AlertMessage
	err := c.BindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// log.Printf("AlertMessage content %+v", notification)

	webhookType := c.Query("type")
	templateName := c.Query("tpl")
	fsURL := c.Query("fsurl")
	if webhookType != "fs" || templateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing parameters"})
		return
	}
	//转告警级别为中文
	AlertLevel := notification.ConvertLevelToInt()
	if notification.Status == "" {

	}
	SendContent, err := json.Marshal(notification)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	// // templateName := notification.GetTemplateName()
	var color, status string
	if strings.Count(string(SendContent), "resolved") > 0 && strings.Count(string(SendContent), "firing") > 0 {
		color = "orange"
	} else if strings.Count(string(SendContent), "resolved") > 0 {
		color = "green"
	} else {
		color = "red"
	}
	if notification.Status == "resolved" {
		status = "恢复"
	} else {
		status = "故障"

	}
	templateFile := filepath.Join("templates", templateName+".tmpl")
	alertTemplate, err := model.NewTemplate(templateFile)
	if err != nil {
		log.Errorf("Failed to load template file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "template loading failed"})
		return
	}

	// for i := range notification.Alerts {
	// 	notification.Alerts[i].Annotations["text"] = messageContent
	// }
	notification.SetTemplate(alertTemplate)
	messageContent, err := notification.Template.Execute(notification)
	fmt.Printf("messageContent:%s\n", messageContent)
	// messageContext, err := alertTemplate.Execute(notification)
	if err != nil {
		log.Errorf("Failed to execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "template execution failed"})
		return
	}
	parsedURL, err := url.Parse(fsURL)
	if err != nil {
		log.Errorf("Failed to parse FeiShu webhook URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fsurl parameter"})
		return
	}

	commonMsg := &model.CommonMessage{
		Platform: config.AppConfig.AlertType,
		Title:    notification.GroupLabels["alertname"],
		Text:     messageContent,
		Level:    AlertLevel,
		Color:    color,
		Status:   status,
	}

	sender := &feishu.FeiShuSender{
		WebhookURL: parsedURL.String(),
	}
	// sender := getSender(commonMsg)
	//fmt.Println(msg)
	// if sender == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported platform"})
	// 	log.Infof("getSender:%s", err)
	// 	return
	// }

	err = sender.SendV2(commonMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successful send alert notification!"})
}

// func getSender(notification *model.CommonMessage) interfaces.MessageSender {
// 	log.Infof("Attempting to get sender for platform: %s", notification.Platform)
// 	switch notification.Platform {
// 	case "dingtalk":
// 		if config.AppConfig.OpenDingding == 1 {
// 			return &dingtalk.DingTalkSender{WebhookURL: config.AppConfig.Dingtalk.WebhookURL}
// 		}
// 	case "feishu":
// 		if config.AppConfig.OpenFeishu == 1 {
// 			return &feishu.FeiShuSender{WebhookURL: config.AppConfig.Feishu.WebhookURL}
// 		}
// 	default:
// 		log.Warnf("Unsupported platform specified: %s", notification.Platform)
// 	}
// 	return nil
// }

// func buildMessageText(notification model.AlertMessage) string {
// 	var buffer strings.Builder
// 	buffer.WriteString(fmt.Sprintf("通知组%s,状态[%s]\n告警项\n\n", notification.GroupKey, notification.Status))
// 	for _, alert := range notification.Alerts {
// 		buffer.WriteString(fmt.Sprintf("摘要：%s\n详情: %s\n", alert.Annotations["summary"], alert.Annotations["description"]))
// 		buffer.WriteString(fmt.Sprintf("开始时间: %s\n\n", alert.StartsAt.Format("2023-12-01 15:04:05")))
// 	}
// 	return buffer.String()
// }

func SendMessageR(message model.AlertMessage, platform, fsurl, phone, email string) {
	AlertMessages := message.Alerts
	var fstext, titleend string
	commonMsg := &model.CommonMessage{
		Platform: platform,
		Title:    titleend,
		Text:     fstext,
	}
	sender := &feishu.FeiShuSender{
		WebhookURL: fsurl,
	}
	err := sender.SendV2(commonMsg)
	if err != nil {
		log.Errorf("发送失败%s", err)
	}
	for _, RMessage := range AlertMessages {
		if RMessage.Status == "resolved" {
			titleend = "故障恢复信息"
		} else {
			titleend = "故障恢复信息"

		}
	}

}
