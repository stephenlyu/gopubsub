package pubsubadapter

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/stephenlyu/gopubsub/message"
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
