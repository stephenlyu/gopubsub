package pubsubadapter

import (
	"github.com/stephenlyu/gopubsub/api"
	"github.com/stephenlyu/gopubsub/redispubsub"
	"github.com/stephenlyu/gopubsub/config"
	"github.com/Sirupsen/logrus"
)

type PubSubAdapter struct {
	server *api.Server
	redisPubSub *redispubsub.RedisPubSub

	translator MessageTranslator
	channels []interface{}
}

func NewPubSubAdapter(redisUrl string, redisPassword string, serverConfig config.Config, serverPort int, serverEndPoint string) *PubSubAdapter {
	server := api.NewServer(serverConfig, serverPort, serverEndPoint)
	redisPubSub := redispubsub.NewRedisPubSub(redisUrl, redisPassword)

	return &PubSubAdapter{
		server: server,
		redisPubSub: redisPubSub,
		translator: &TransparentTranslator{},
	}
}

func (this *PubSubAdapter) SetChannels(channels ...interface{}) {
	this.channels = channels
}

func (this *PubSubAdapter) SetTranslator(translator MessageTranslator) {
	this.translator = translator
}

func (this *PubSubAdapter) OnMessage(channel string, data []byte) {
	m := this.translator.TranslateTo(channel, data)
	if m == nil {
		return
	}

	this.server.Publish(*m)
}

func (this *PubSubAdapter) redisSubscribeLoop() {
	for {
		err := this.redisPubSub.Subscribe(this, this.channels...)
		if err != nil {
			logrus.Errorf("PubSubAdapter.redisSubscribeLoop error: %s", err.Error())
		}
	}
}

func (this *PubSubAdapter) Start() {
	go this.redisSubscribeLoop()
	this.server.Start()
}
