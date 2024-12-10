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
	FeishuRobots  []RobotConfig    `yaml:"feishu_robots"`
	Templates     []TemplateConfig `yaml:"templates"`
	DB            DBConfig         `yaml:"db"`
}

type DingtalkConfig struct {
	WebhookURL string `yaml:"dd_url"`
}

type FeishuConfig struct {
	WebhookURL string `yaml:"fs_url"`
}
type RobotConfig struct {
	Name       string `yaml:"name"`
	WebhookURL string `yaml:"webhook_url"`
}

type TemplateConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
type DBConfig struct {
	Driver          string `yaml:"driver"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
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
