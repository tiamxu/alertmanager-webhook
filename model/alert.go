package model

import "time"

// prometheus告警类型
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
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
