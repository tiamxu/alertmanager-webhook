package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/alertmanager-webhook/service"
)

var alertService = service.NewAlertService()

func PrometheusAlert(c *gin.Context) {
	var notification model.AlertMessage
	if err := c.BindJSON(&notification); err != nil {
		handleError(c, http.StatusBadRequest, "failed to parse JSON", err)
		return
	}

	webhookType := c.Query("type")
	templateName := c.Query("tpl")
	fsURL := c.Query("fsurl")
	atSomeOne := c.Query("at")
	split := c.Query("split")

	messageData, err := alertService.ProcessAlert(&notification, webhookType, templateName, fsURL, atSomeOne, split)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "processing alert failed", err)
		return
	}

	response := model.Response{
		Code: http.StatusOK,
		Msg:  "successful send alert notification!",
		Data: messageData,
	}
	c.JSON(http.StatusOK, response)
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Errorf("Error %s: %v", message, err)
	response := model.Response{
		Code: statusCode,
		Msg:  fmt.Sprintf("%s: %v", message, err),
		Data: nil,
	}
	c.JSON(statusCode, response)
}
