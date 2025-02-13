package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/alertmanager-webhook/service"
	"github.com/tiamxu/kit/log"
)

type AlertHandler struct {
	alertService service.AlertService
}

func NewAlertHandler(service service.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: service,
	}
}

func (h *AlertHandler) PrometheusAlert(c *gin.Context) {
	var notification model.AlertMessage
	if err := c.BindJSON(&notification); err != nil {
		handleError(c, http.StatusBadRequest, "无效的告警数据格式", err)
		return
	}

	// 参数验证
	webhookType := c.Query("type")
	if webhookType == "" {
		handleError(c, http.StatusBadRequest, "缺少必要参数: type", nil)
		return
	}

	templateName := c.Query("tpl")

	// 根据不同的 webhookType 获取对应的 URL
	var webhookURL string
	switch webhookType {
	case "fs":
		webhookURL = c.Query("fsurl")
	case "dd":
		webhookURL = c.Query("ddurl")
	default:
		handleError(c, http.StatusBadRequest, "不支持的告警类型", nil)
		return
	}

	if webhookURL == "" {
		handleError(c, http.StatusBadRequest, "缺少必要的 Webhook URL 参数", nil)
		return
	}

	atSomeOne := c.Query("at")
	split := c.Query("split")

	messageData, err := h.alertService.ProcessAlert(&notification, webhookType, templateName, webhookURL, atSomeOne, split)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "处理告警失败", err)
		return
	}

	response := model.Response{
		Code: http.StatusOK,
		Msg:  "告警通知发送成功！",
		Data: messageData,
	}
	c.JSON(http.StatusOK, response)
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
		log.Errorf("处理告警错误: %s, 详细信息: %v", message, err)
	}

	response := model.Response{
		Code: statusCode,
		Msg:  errMsg,
		Data: nil,
	}
	c.JSON(statusCode, response)
}
