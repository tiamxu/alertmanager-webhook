package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env           string           `yaml:"env"`
	LogLevel      string           `yaml:"log_level"`
	ListenAddress string           `yaml:"listen_address"`
	HotReload     bool             `yaml:"hot_reload"`
	AlertType     string           `yaml:"alert_type"`
	OpenDingding  int              `yaml:"open_dingding"`
	OpenFeishu    int              `yaml:"open_feishu"`
	Dingtalk      DingtalkConfig   `yaml:"dingtalk"`
	Feishu        FeishuConfig     `yaml:"feishu"`
	Templates     []TemplateConfig `yaml:"templates"`
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

var AppConfig *Config

func Load() error {
	filename := "./config/config.yaml"
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	AppConfig = &Config{}
	err = yaml.Unmarshal(data, AppConfig)
	if err != nil {
		return err
	}

	return nil
}
