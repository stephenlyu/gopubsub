
package main

import (
	"github.com/stephenlyu/gopubsub/api"
	"github.com/stephenlyu/gopubsub/config"
	"github.com/Sirupsen/logrus"
	"github.com/stephenlyu/gopubsub/message"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "20060102 15:04:05.999999999"})

	server := api.NewServer(*config.DEFAULT_CONFIG, 8080, "/quote")
	go server.Start()

	counter := 0
	for {
		server.Publish(message.Message{Subject: "a", Data:counter})
		counter++
		//time.Sleep(time.Second)
		//server.Publish(message.Message{Subject: "b", Data:counter})
		//counter++
		time.Sleep(time.Millisecond)
	}
}
