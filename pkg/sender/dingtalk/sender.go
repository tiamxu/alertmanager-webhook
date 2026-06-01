package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/kit/log"
)

// 全局 HTTP Client，复用连接
var httpClient = &http.Client{Timeout: 10 * time.Second}

// BotConfig 机器人配置
type BotConfig struct {
	Secret string
}

// botsSecret 读写锁，保护并发安全
var (
	botsSecretMu sync.RWMutex
	botsSecret   = make(map[string]BotConfig)
)

// SetBotsConfig 设置机器人配置（由 main.go 调用）
func SetBotsConfig(bots map[string]BotConfig) {
	botsSecretMu.Lock()
	defer botsSecretMu.Unlock()
	botsSecret = bots
}

// ValidateBot 验证 bot 是否有效
func ValidateBot(bot string) bool {
	botsSecretMu.RLock()
	defer botsSecretMu.RUnlock()
	_, exists := botsSecret[bot]
	return exists
}

// DingTalkSender 钉钉消息发送器
type DingTalkSender struct {
	WebhookURL string
	Bot        string
}

// NewDingtalkSender 创建钉钉发送器
func NewDingtalkSender(webhookURL, bot string) *DingTalkSender {
	return &DingTalkSender{WebhookURL: webhookURL, Bot: bot}
}

// Send 发送消息（主入口）
func (d *DingTalkSender) Send(message *model.CommonMessage) error {
	if message.Platform != "dingtalk" {
		return fmt.Errorf("invalid platform for DingTalkSender")
	}

	msg, err := d.buildMessage(message)
	if err != nil {
		return fmt.Errorf("构建消息失败: %v", err)
	}

	secret, ok := GetBotSecret(d.Bot)
	if !ok {
		return fmt.Errorf("failed to get secret for bot: %s", d.Bot)
	}
	webhookURL, err := d.generateSignedURL(secret)
	if err != nil {
		return fmt.Errorf("生成签名URL失败: %v", err)
	}
	return d.sendRequest(webhookURL, msg)
}

// buildMessage 构建钉钉消息
func (d *DingTalkSender) buildMessage(message *model.CommonMessage) (*model.DDMessage, error) {
	title := fmt.Sprintf("%s [%s] [%s]", message.Title, message.Level, message.Status)
	titleend := fmt.Sprintf("### <font size=6 color='%s'>%s [%s] [%s]</font>\n\n",
		message.Color, message.Title, message.Level, message.Status)

	content, atMobiles, isAtAll := d.processAtContent(message.Text, message.AtSomeOne)
	content = titleend + content

	return &model.DDMessage{
		MsgType: "markdown",
		Markdown: model.Md{
			Title: title,
			Text:  content,
		},
		At: model.At{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}, nil
}

// processAtContent 处理@提醒和消息内容
func (d *DingTalkSender) processAtContent(content, atSomeOne string) (string, []string, bool) {
	if atSomeOne == "" {
		return content, []string{"18888888888"}, true
	}

	atMobiles := strings.Split(atSomeOne, ",")
	var atText strings.Builder
	for _, phone := range atMobiles {
		atText.WriteString(" @")
		atText.WriteString(phone)
	}
	atText.WriteString(".")
	finalContent := content + atText.String()

	return finalContent, atMobiles, false
}

// generateSignedURL 生成带签名的URL
func (d *DingTalkSender) generateSignedURL(secret string) (string, error) {
	timestamp := time.Now().UnixMilli()
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(stringToSign)); err != nil {
		return "", fmt.Errorf("failed to write hmac: %v", err)
	}

	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	signEscaped := url.QueryEscape(sign)

	return fmt.Sprintf("%s&timestamp=%d&sign=%s", d.WebhookURL, timestamp, signEscaped), nil
}

// sendRequest 发送请求
func (d *DingTalkSender) sendRequest(webhookURL string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.WithFields(log.Fields{
			"module": "dingtalk",
			"action": "send",
			"error":  err.Error(),
		}).Error("发送消息失败")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"module": "dingtalk",
			"action": "send",
			"error":  err.Error(),
		}).Error("读取响应失败")
		return err
	}

	log.WithFields(log.Fields{
		"module":   "dingtalk",
		"action":   "send",
		"response": string(body),
	}).Info("消息发送成功")
	return nil
}

// GetPlatform 获取平台标识
func (d *DingTalkSender) GetPlatform() string {
	return "dingtalk"
}

// GetBotSecret 获取指定 bot 的 secret
func GetBotSecret(bot string) (string, bool) {
	botsSecretMu.RLock()
	defer botsSecretMu.RUnlock()
	config, exists := botsSecret[bot]
	if !exists || config.Secret == "" {
		return "", false
	}
	return config.Secret, true
}
