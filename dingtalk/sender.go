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
	"strconv"
	"strings"
	"time"

	"github.com/tiamxu/alertmanager-webhook/model"
	"github.com/tiamxu/kit/log"
)

type DingTalkSender struct {
	WebhookURL string
	Secret     string
}

const default_DingTalk_Secret = "SEC0488bbb01d1bdb222619fce742061687a1e591d5a914923781495bde1128c8bf"

// 签名
func (d *DingTalkSender) sign(timestamp int64) string {
	d.Secret = default_DingTalk_Secret
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.Secret)
	h := hmac.New(sha256.New, []byte(d.Secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (d *DingTalkSender) Send(message *model.CommonMessage) error {
	if message.Platform != "dingtalk" {
		return fmt.Errorf("invalid platform for DingTalkSender")
	}
	// 添加title的样式
	title := fmt.Sprintf("%s [%s] [%s]", message.Title, message.Level, message.Status)

	titleend := fmt.Sprintf("### <font size=6 color='%s'>%s [%s] [%s]</font>\n\n", message.Color, message.Title, message.Level, message.Status)
	SendContent := titleend + message.Text

	Atall := true
	atMobile := []string{"18888888888"}

	AtSomeOne := message.AtSomeOne
	if AtSomeOne != "" {
		atMobile = strings.Split(AtSomeOne, ",")
		AtText := ""
		for _, phoneN := range atMobile {
			AtText += " @" + phoneN
		}
		SendContent = SendContent + AtText + "."
		Atall = false
	}
	var msg interface{}

	// 构建钉钉消息
	// msg = map[string]interface{}{
	// 	"msgtype": "markdown",
	// 	"markdown": map[string]string{
	// 		"title": message.Title,
	// 		"text":  d.formatMessage(message),
	// 	},
	// 	"at": map[string]interface{}{
	// 		"atMobiles": d.parseAtList(message.AtSomeOne),
	// 		"isAtAll":   false,
	// 	},
	// }
	msg = model.DDMessage{
		MsgType: "markdown",
		Markdown: model.Md{
			Title: title,
			Text:  SendContent,
		},
		At: model.At{
			AtMobiles: atMobile,
			IsAtAll:   Atall,
		},
	}
	// 添加签名
	timestamp := time.Now().UnixNano() / 1e6
	sign := d.sign(timestamp)
	webhookURL := d.addSignature(timestamp, sign)

	// 发送消息
	return d.sendRequest(webhookURL, msg)
}

func (d *DingTalkSender) addSignature(timestamp int64, sign string) string {
	baseURL, _ := url.Parse(d.WebhookURL)
	query := baseURL.Query()
	query.Set("timestamp", strconv.FormatInt(timestamp, 10))
	query.Set("sign", sign)
	baseURL.RawQuery = query.Encode()
	return baseURL.String()
}

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
