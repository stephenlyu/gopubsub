package connmanager

import (
	"golang.org/x/net/websocket"
	"github.com/stephenlyu/gopubsub/message"
	"strings"
	"github.com/Sirupsen/logrus"
)

type Connection struct {
	manager *ConnManager
	Id     string
	Socket *websocket.Conn
	SendCh   chan *message.Message
}

func (this *Connection) handleMessage(msg *message.Message) {
	switch msg.Subject {
	case "PING":
		this.manager.Send(message.Message{Subject: "PONG", Data: msg.Data})
	case "SUBSCRIBE":
		if subjectStr, ok := msg.Data.(string); ok {
			subjects := strings.Split(subjectStr, ",")
			this.manager.Subscribe(this, subjects)
			this.SendCh <- &message.Message{Subject:"SUBSCRIBE", Data: "OK"}
		} else {
			this.SendCh <- &message.Message{Subject:"SUBSCRIBE", Data: "BAD SUBJECTS"}
		}
	case "UNSUBSCRIBE":
		if subjectStr, ok := msg.Data.(string); ok {
			subjects := strings.Split(subjectStr, ",")
			this.manager.UnSubscribe(this, subjects)
			this.SendCh <- &message.Message{Subject:"UNSUBSCRIBE", Data: "OK"}
		} else {
			this.SendCh <- &message.Message{Subject:"UNSUBSCRIBE", Data: "BAD SUBJECTS"}
		}
	}
}

func (this *Connection) Read() {
	defer func() {
		this.manager.unregisterCh <- this
		this.Socket.Close()
	}()

	var message *message.Message
	for {
		err := websocket.JSON.Receive(this.Socket, &message)
		if err != nil {
			logrus.Errorf("Connection.Read error: %+v", err)
			break
		}
		this.handleMessage(message)
	}
}

func (this *Connection) Write() {
	defer func() {
		this.Socket.Close()
	}()

	for {
		select {
		case message, ok := <- this.SendCh:
			if !ok {
				logrus.Error("Connection.Write SendCh closed.")
				return
			}

			// TODO: 支持重试
			err := websocket.JSON.Send(this.Socket, message)
			if err != nil {
				logrus.Errorf("Connection.Write error: %+v", err)
				return
			}
		}
	}
}

func (this *Connection) Run() {
	go this.Write()
	this.Read()
}
