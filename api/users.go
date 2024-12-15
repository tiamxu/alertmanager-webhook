package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/tiamxu/alertmanager-webhook/service"
)

type TenantAccessMeg struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}
type TenantAccessResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

const APP_ID string = "cli_a1171daf3839d00e"
const APP_SECRET string = "EJkjsUXjk5soHByhuL1ftfzKIWPM5wdn"

func GetUserIDsByAttributes(c *gin.Context) {
	var userAttributes service.UserAttributes
	userIdType := c.Query("user_id_type")
	err := c.BindJSON(&userAttributes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := lark.NewClient(APP_ID, APP_SECRET)

	response, err := userAttributes.UserIDsByAttributes(client, userIdType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userIDs := response.Data.UserList

	c.JSON(http.StatusOK, gin.H{"user_ids": userIDs})

}
func GetUserIDsByDepartment(c *gin.Context) {
	departmentId := c.Query("department_id")
	if departmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing department_id"})
		return
	}

	request := &service.ListUserByDepartmentRequest{
		DepartmentId: departmentId,
	}
	client := lark.NewClient(APP_ID, APP_SECRET)
	response, err := service.ListUserByDepartment(client, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userIDs := make([]string, 0)
	for _, user := range response.FindByDepartmentUserResponse {
		userIDs = append(userIDs, *user.OpenId) // 使用 OpenId 或 UnionId 根据需求选择
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": userIDs})
}

func GetTenantAccessTokenBySelfBuiltApp() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetTenantAccessTokenBySelfBuiltApp(context.Background(), &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}
