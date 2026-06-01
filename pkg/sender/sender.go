package sender

import "github.com/tiamxu/alertmanager-webhook/model"

// MessageSender 消息发送器接口
type MessageSender interface {
	Send(message *model.CommonMessage) error
	GetPlatform() string
}