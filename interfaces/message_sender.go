package interfaces

import "github.com/tiamxu/alertmanager-webhook/model"

type MessageSender interface {
	Send(message *model.CommonMessage) error
}
