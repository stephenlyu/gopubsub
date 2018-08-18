package main

import (
	"fmt"
	"github.com/stephenlyu/gopubsub/api"
	"github.com/stephenlyu/gopubsub/message"
	"github.com/Sirupsen/logrus"
)

type _callback struct {
}

func (this _callback) OnMessage(message *message.Message, raw []byte) {
	fmt.Printf("OnMessage: %+v\n", message)
}

func (this _callback) OnError(err error) {
	fmt.Println(err)
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "20060102 15:04:05.999999999"})

	client := api.NewClient("localhost", api.DEFAULT_PORT, api.DEFAULT_END_POINT, &_callback{})
	err := client.Start()
	if err != nil {
		panic(err)
	}

	client.Subscribe([]string{"a", "b"})

	ch := make(chan struct{})
	<- ch
}
