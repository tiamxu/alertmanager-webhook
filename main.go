package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/api"
	"github.com/tiamxu/alertmanager-webhook/service"
	httpkit "github.com/tiamxu/kit/http"
	"github.com/tiamxu/kit/log"
)

var cfg *Config

func init() {
	loadConfig()
	if err := cfg.Initial(); err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}
}
func main() {

	// 初始化 service 和 handler
	alertService := service.NewAlertService()          // 创建 service 实例
	alertHandler := api.NewAlertHandler(*alertService) // 创建 handler 实例
	router := httpkit.NewGin(cfg.HttpSrv)

	// 路由注册
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 使用 handler 的方法替代直接调用
	router.POST("/webhook", alertHandler.PrometheusAlert)

	router.POST("/users/get_id", api.GetUserIDsByAttributes)
	router.GET("/get_user_ids", api.GetUserIDsByDepartment)

	// 启动服务器
	srv := httpkit.StartServer(router, cfg.HttpSrv)
	log.Infoln("Server listen: ", cfg.HttpSrv.Address)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infoln("Shutting down server...")
	httpkit.ShutdownServer(srv)
	log.Infoln("Server exited")

}
