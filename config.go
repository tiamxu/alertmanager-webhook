package main

import (
	"fmt"
	"os"

	httpkit "github.com/tiamxu/kit/http"
	"github.com/tiamxu/kit/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env          string                  `yaml:"env"`
	LogLevel     string                  `yaml:"log_level"`
	HttpSrv      httpkit.GinServerConfig `yaml:"http_srv"`
	HotReload    bool                    `yaml:"hot_reload"`
	AlertType    string                  `yaml:"alert_type"`
	OpenDingding int                     `yaml:"open_dingding"`
	OpenFeishu   int                     `yaml:"open_feishu"`
	Dingtalk     DingtalkConfig          `yaml:"dingtalk"`
	Feishu       FeishuConfig            `yaml:"feishu"`
	Templates    []TemplateConfig        `yaml:"templates"`
}

type DingtalkConfig struct {
	WebhookURL string `yaml:"dd_url"`
}

type FeishuConfig struct {
	WebhookURL string `yaml:"fs_url"`
}

type TemplateConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

var configPath = "config/config.yaml"

func (c *Config) Initial() (err error) {

	defer func() {
		if err == nil {
			log.Printf("config initialed, env: %s", cfg.Env)
		}
	}()
	//日志
	// if level, err := logrus.ParseLevel(c.LogLevel); err != nil {
	// 	return err
	// } else {
	// 	log.DefaultLogger().SetLevel(level)
	// }
	err = log.InitLogger(&log.Config{
		Level:      "info",
		Type:       "stdout", // "file" 或 "stdout"
		Format:     "json",   // "json" 或 "text"
		FilePath:   "logs",
		FileName:   "alert.log",
		MaxSize:    100, // 每个文件最大 100MB
		MaxAge:     7,   // 保留7天
		MaxBackups: 10,  // 保留10个备份
		Compress:   true,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log.SetGlobalFields(log.Fields{
		"appname": "alert",
	})

	return nil
}
func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
