package sender

import (
	"fmt"

	"github.com/tiamxu/alertmanager-webhook/pkg/sender/dingtalk"
	"github.com/tiamxu/alertmanager-webhook/pkg/sender/feishu"
)

// NewSender 根据类型创建对应的发送器
func NewSender(senderType, webhookURL, bot string) (MessageSender, error) {
	switch senderType {
	case "fs":
		return feishu.NewFeishuSender(webhookURL), nil
	case "dd":
		if !dingtalk.ValidateBot(bot) {
			return nil, fmt.Errorf("invalid bot: %s", bot)
		}
		return dingtalk.NewDingtalkSender(webhookURL, bot), nil
	default:
		return nil, fmt.Errorf("unsupported sender type: %s", senderType)
	}
}