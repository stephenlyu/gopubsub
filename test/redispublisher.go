package main

import (
	"github.com/stephenlyu/gopubsub/redispubsub"
	"fmt"
	"time"
)

func main() {
	redisPubSub := redispubsub.NewRedisPubSub("", "")

	counter := 0
	for {
		redisPubSub.Publish("a", []byte(fmt.Sprintf("%d", counter)))
		counter++
		redisPubSub.Publish("b", []byte(fmt.Sprintf("%d", counter)))
		counter++
		time.Sleep(time.Millisecond)
	}
}
