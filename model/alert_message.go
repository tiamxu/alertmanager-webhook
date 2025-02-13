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

type Alerts []Alert

// 定义告警级别的优先级映射
var levelPriority = map[string]int{
	"信息":   0,
	"警告":   1,
	"一般严重": 2,
	"严重":   3,
	"灾难":   4,
}

// 获取告警级别的优先级值，默认返回最小优先级
func getLevelPriority(alert Alert) int {
	levelStr, exists := alert.Labels["level"]
	if !exists {
		return 0 // 如果没有找到 level 字段，返回默认值
	}

	nLevel, err := strconv.Atoi(levelStr)
	if err != nil || nLevel < 0 {
		return 0 // 如果转换失败或数值无效，返回默认值
	}

	return nLevel
}

// 实现 sort.Interface 接口的方法
func (a Alerts) Len() int {
	return len(a)
}

func (a Alerts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Alerts) Less(i, j int) bool {
	// 获取两个告警的优先级
	priorityI := getLevelPriority(a[i])
	priorityJ := getLevelPriority(a[j])

	// 从大到小排序，所以这里是比较priorityJ是否小于priorityI
	return priorityJ < priorityI
}
