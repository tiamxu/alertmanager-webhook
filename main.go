package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	api "github.com/tiamxu/alertmanager-webhook/api"
	"github.com/tiamxu/alertmanager-webhook/config"
	"github.com/tiamxu/alertmanager-webhook/log"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化日志级别
	if level, err := logrus.ParseLevel(config.AppConfig.LogLevel); err == nil {
		log.DefaultLogger().SetLevel(level)
	} else {
		log.Fatalf("无效的日志级别设置: %v", err)
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.POST("/webhook", api.HandlerWebhook)

	srv := &http.Server{
		Addr:    config.AppConfig.ListenAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe failed: %v", err)
		}
	}()
	log.Infof("服务正在监听地址: %s", config.AppConfig.ListenAddress)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("收到终止信号,开始停止HTTP服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	}

	log.Infoln("HTTP服务器已成功停止")
}
