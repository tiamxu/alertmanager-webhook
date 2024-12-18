package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/log"
	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/alertmanager-webhook/pkg/e"
	"github.com/tiamxu/alertmanager-webhook/service"
)

var alertService = service.NewAlertService()

func PrometheusAlert(c *gin.Context) {
	webhookType := c.Query("type")
	templateName := c.Query("tpl")
	fsURL := c.Query("fsurl")
	atSomeOne := c.Query("at")
	split := c.Query("split")

	var notification model.AlertMessage

	if err := c.BindJSON(&notification); err != nil {
		// handleError(c, http.StatusBadRequest, "failed to parse JSON", err)
		code := e.ERROR
		// return nil, fmt.Errorf("invalid or missing parameters")
		res := model.Response{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "invalid or missing parameters",
		}
		c.JSON(http.StatusBadRequest, res)

	} else {
		response, _ := alertService.ProcessAlert(&notification, webhookType, templateName, fsURL, atSomeOne, split)

		c.JSON(http.StatusOK, response)

	}

	// if err != nil {
	// 	handleError(c, http.StatusInternalServerError, "processing alert failed", err)
	// 	return
	// }

	// templateName := notification.GetTemplateName()

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
