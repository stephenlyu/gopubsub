package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/stephenlyu/gopubsub/adapter"
	"github.com/stephenlyu/gopubsub/config"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "20060102 15:04:05.999999999"})
	adapter := pubsubadapter.NewPubSubAdapter("", "", *config.DEFAULT_CONFIG, 0, "")
	adapter.SetChannels("a", "b")

	adapter.Start()
}
