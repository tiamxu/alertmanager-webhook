package model

type Config struct {
	ENV           string `yaml:"env"`
	LogLevel      string `yaml:"log_level"`
	ListenAddress string `yaml:"listen_address"`
	HotReload     bool   `yaml:"hot_reload"`
	AlertType     string `yaml:"alert_type"`
	OpenDingding  string `yaml:"open_dingding"`
	OpenFeishu    string `yaml:"open_feishu"`
	Dingding      `yaml:"dingding"`
	Feishu        `yaml:"feishu"`
}
type Dingding struct {
	DdUrl string `yaml:"dd_url"`
}
type Feishu struct {
	FsUrl string `yaml:"fs_url"`
}
