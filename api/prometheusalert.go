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
	if err := c.BindJSON(&notification); err != nil {
		handleError(c, http.StatusBadRequest, "failed to parse JSON", err)
		return
	}
	// log.Printf("AlertMessage content %+v", notification)

	webhookType := c.Query("type")
	templateName := c.Query("tpl")
	fsURL := c.Query("fsurl")
	atSomeOne := c.Query("at")
	split := c.Query("split")

	if webhookType != "fs" || templateName == "" {
		handleError(c, http.StatusBadRequest, "invalid or missing parameters", fmt.Errorf("接口参数异常"))
		return
	}

	// 转告警级别为中文并设置消息颜色和状态
	level := notification.ConvertLevelToInt()
	color, status := getAlertColorAndStatus(notification)

	// templateName := notification.GetTemplateName()
	//加载模版文件
	templateFile := filepath.Join("templates", templateName+".tmpl")
	alertTemplate, err := model.NewTemplate(templateFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "template loading failed", err)
		return
	}
	notification.SetTemplate(alertTemplate)
	// fmt.Printf("messageContent:%s\n", messageContent)
	// messageContext, err := alertTemplate.Execute(notification)

	// 解析 FeiShu webhook URL
	parsedURL, err := url.Parse(fsURL)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid fsurl parameter", err)
		return
	}

	// 构建和发送消息
	if split == "true" {
		for _, alert := range notification.Alerts {
			// singleAlertMsg := notification
			notification.Alerts = []model.Alert{alert}
			notification.Status = alert.Status
			color, status := getAlertColorAndStatus(notification)

			sendAlertMessage(c, parsedURL.String(), level, color, status, atSomeOne, alert.Labels["alertname"], notification, alertTemplate)
		}
	} else {
		sendAlertMessage(c, parsedURL.String(), level, color, status, atSomeOne, notification.GroupLabels["alertname"], notification, alertTemplate)
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

// 构建消息内容
func sendAlertMessage(c *gin.Context, fsurl, level, color, status, atSomeOne, title string, message interface{}, tmpl *model.Template) {
	messageContent, err := tmpl.Execute(message)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "template execution failed", err)
		return
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
		handleError(c, http.StatusInternalServerError, "message sending failed", err)
	}
}

func getAlertColorAndStatus(notification model.AlertMessage) (string, string) {
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

func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Errorf("Error %s: %v", message, err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprintf("%s: %v", message, err)})
}
