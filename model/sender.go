package model

// MessageSender 消息发送器接口
// 注意：此接口已迁移至 pkg/sender.MessageSender，此处保留仅为向后兼容
type MessageSender interface {
	Send(message *CommonMessage) error
	GetPlatform() string
}
