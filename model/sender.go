package model

type MessageSender interface {
	Send(message *CommonMessage) error
	GetPlatform() string
}
