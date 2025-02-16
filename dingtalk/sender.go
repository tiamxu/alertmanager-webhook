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
	"time"

	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/kit/log"
)

// 定义钉钉机器人secret
var botsSecret = map[string]string{
	"bot1": "SEC0488bbb01d1bdb222619fce742061687a1e591d5a914923781495bde1128c8bf",
	"bot2": "",
}

type DingTalkSender struct {
	WebhookURL string
	Bot        string
}

// ValidateBot 验证 bot 是否有效
func ValidateBot(bot string) bool {
	secret, exists := botsSecret[bot]
	return exists && secret != ""
}

// GetBotSecret 获取指定 bot 的 secret
func GetBotSecret(bot string) (string, bool) {
	secret, exists := botsSecret[bot]
	if !exists || secret == "" {
		return "", false
	}
	return secret, true
}

// const default_DingTalk_Secret = "SEC0488bbb01d1bdb222619fce742061687a1e591d5a914923781495bde1128c8bf"

func (d *DingTalkSender) Send(message *model.CommonMessage) error {
	if message.Platform != "dingtalk" {
		return fmt.Errorf("invalid platform for DingTalkSender")
	}

	// 构建消息内容
	msg, err := d.buildMessage(message)
	if err != nil {
		return fmt.Errorf("构建消息失败: %v", err)
	}

	// 获取 secret
	secret, ok := GetBotSecret(d.Bot)
	if !ok {
		return fmt.Errorf("failed to get secret for bot: %s", d.Bot)
	}
	// 添加签名
	// timestamp := time.Now().UnixMilli()
	// sign, err := d.generateSign(secret, timestamp)
	// if err != nil {
	// 	return fmt.Errorf("failed to generate sign: %v", err)
	// }
	// webhookURL := fmt.Sprintf("%s&timestamp=%d&sign=%s", d.WebhookURL, timestamp, sign)
	// 生成签名URL
	webhookURL, err := d.generateSignedURL(secret)
	if err != nil {
		return fmt.Errorf("生成签名URL失败: %v", err)
	}
	// 发送消息
	return d.sendRequest(webhookURL, msg)
}

// buildMessage 构建钉钉消息
func (d *DingTalkSender) buildMessage(message *model.CommonMessage) (*model.DDMessage, error) {
	// 构建标题
	title := fmt.Sprintf("%s [%s] [%s]", message.Title, message.Level, message.Status)
	titleend := fmt.Sprintf("### <font size=6 color='%s'>%s [%s] [%s]</font>\n\n",
		message.Color, message.Title, message.Level, message.Status)

	// 处理@人员和消息内容
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

	// 分割手机号
	atMobiles := strings.Split(atSomeOne, ",")

	// 构建@文本
	var atText strings.Builder
	for _, phone := range atMobiles {
		atText.WriteString(" @")
		atText.WriteString(phone)
	}
	atText.WriteString(".")

	// 将@文本添加到内容末尾
	finalContent := content + atText.String()

	return finalContent, atMobiles, false
}

// generateSign 生成签名
func (d *DingTalkSender) generateSign(secret string, timestamp int64) (string, error) {

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(stringToSign)); err != nil {
		return "", fmt.Errorf("failed to write hmac: %v", err)
	}
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(sign), nil
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

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Errorln("[dingtalk]", err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("[dingtalk]", err.Error())
		return err
	}

	log.Infoln("[dingtalk]", string(body))
	return nil
}
func (d *DingTalkSender) GetPlatform() string {
	return "dingtalk"
}
