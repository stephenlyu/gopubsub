package pubsubadapter

import (
	"github.com/stephenlyu/gopubsub/message"
	"github.com/Sirupsen/logrus"
	"encoding/json"
)

type MessageTranslator interface {
	TranslateTo(channel string, data []byte) *message.Message
}

type TransparentTranslator struct {
}

func (this TransparentTranslator) TranslateTo(channel string, data []byte) *message.Message {
	var r interface{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		logrus.Errorf("PubSubAdapter.OnMessage error: %s", err.Error())
		return nil
	}

	return &message.Message{Subject: channel, Data: r}
}
