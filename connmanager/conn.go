package connmanager

import (
	"golang.org/x/net/websocket"
	"github.com/stephenlyu/gopubsub/message"
	"strings"
	"github.com/Sirupsen/logrus"
	"time"
	"encoding/json"
)

type Connection struct {
	manager *ConnManager
	Id     string
	Socket *websocket.Conn
	SendCh   chan interface{}
}

func (this *Connection) handleMessage(msg *message.Message) {
	switch msg.Subject {
	case "PING":
		logrus.Infof("Connection.handleMessage PING(%s)", this.Id)
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
		this.Socket.Close()
		this.manager.unregisterCh <- this
	}()

	var message *message.Message
	for {
		err := websocket.JSON.Receive(this.Socket, &message)
		if err != nil {
			logrus.Errorf("Connection.Read conn %s error: %+v", this.Id, err)
			break
		}
		this.handleMessage(message)
	}
}

func (this *Connection) Write() {
	defer func() {
		this.Socket.Close()
		this.manager.unregisterCh <- this
	}()

	var err error
	for {
		select {
		case message, ok := <- this.SendCh:
			if !ok {
				logrus.Errorf("Connection.Write conn %s SendCh closed.", this.Id)
				return
			}
			logrus.Infof("Connection.Write conn %s writing..", this.Id)
			startT := time.Now()
			var data []byte
			this.Socket.SetWriteDeadline(time.Now().Add(time.Second * 10))
			switch message.(type) {
			case []byte:
				data = message.([]byte)
			case string:
				data = []byte(message.(string))
			default:
				data, _ = json.Marshal(message)
				if this.manager.config.SupportZip {
					data, _ = GzipEncode(data)
				}
			}
			err = websocket.Message.Send(this.Socket, data)

			// TODO: 支持重试
			if err != nil {
				logrus.Errorf("Connection.Write conn %s error: %+v", this.Id, err)
				return
			}
			logrus.Infof("Connection.Write conn %s timecost: %s", this.Id, time.Now().Sub(startT))
		}
	}
}

func (this *Connection) Run() {
	go this.Write()
	this.Read()
}
