package model

import (
	"strconv"
)

// prometheus告警类型
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type AlertMessage struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Template          *Template
}

func (n *AlertMessage) GetTemplateName() string {
	if val, ok := n.CommonAnnotations["template"]; ok {
		return val
	}
	return "default"
}

func (n *AlertMessage) SetTemplate(tmpl *Template) {
	n.Template = tmpl
}

func (n *AlertMessage) GetFeishuRobotName() string {
	if val, ok := n.CommonLabels["robot"]; ok {
		return val
	}
	return "robot1"
}

var AlertLevel = map[int]string{
	0: "信息",
	1: "警告",
	2: "一般严重",
	3: "严重",
	4: "灾难",
}

// 关于告警级别level共有5个级别,0-4,0 信息,1 警告,2 一般严重,3 严重,4 灾难
func (n *AlertMessage) ConvertLevelToInt() string {
	levelStr, exists := n.CommonLabels["level"]
	if !exists {
		return "未知" // 如果没有找到 level 字段，返回默认值
	}

	nLevel, err := strconv.Atoi(levelStr)
	if err != nil || nLevel < 0 || nLevel > 4 {
		return "未知" // 如果转换失败或超出范围，返回默认值
	}

	return AlertLevel[nLevel]
}
