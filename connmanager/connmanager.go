package connmanager

import (
	"encoding/json"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	. "github.com/stephenlyu/gopubsub/config"
	. "github.com/stephenlyu/gopubsub/message"
	"golang.org/x/net/websocket"
)

type SubUnSub struct {
	subjects []string
	Conn     *Connection
}

type ConnManager struct {
	config Config

	connections map[string]*Connection

	subjectSubscribers map[string]map[*Connection]bool // key: connection id 	value: map from connection pointer to bool
	connectionSubjects map[string]map[string]bool      // key: connection id 	value: map from subject to bool

	registerCh   chan *Connection
	unregisterCh chan *Connection

	subscribeCh   chan *SubUnSub
	unsubscribeCh chan *SubUnSub

	messageCh chan *Message
}

func NewConnManager(config Config) *ConnManager {
	return &ConnManager{
		config: config,

		connections: make(map[string]*Connection),

		subjectSubscribers: make(map[string]map[*Connection]bool),
		connectionSubjects: make(map[string]map[string]bool),

		registerCh:   make(chan *Connection),
		unregisterCh: make(chan *Connection),

		subscribeCh:   make(chan *SubUnSub),
		unsubscribeCh: make(chan *SubUnSub),

		messageCh: make(chan *Message, config.ConnectionSendBufSize),
	}
}

func (this *ConnManager) removeConnectionFromSubjectSubscribers(subject string, conn *Connection) {
	if m, ok := this.subjectSubscribers[subject]; ok {
		delete(m, conn)
	}
}

func (this *ConnManager) Start() {
	clearConnection := func(conn *Connection) {
		delete(this.connections, conn.Id)

		// Remove connection from subject subscribers
		if subjects, ok := this.connectionSubjects[conn.Id]; ok {
			for subject := range subjects {
				this.removeConnectionFromSubjectSubscribers(subject, conn)
				delete(this.connectionSubjects[conn.Id], subject)
			}
		}
	}

	for {
		select {
		case conn := <-this.registerCh:
			logrus.Infof("ConnManager conn %s connected.", conn.Id)
			this.connections[conn.Id] = conn
			this.connectionSubjects[conn.Id] = make(map[string]bool)

		case conn := <-this.unregisterCh:
			if _, ok := this.connections[conn.Id]; ok {
				close(conn.SendCh)
				clearConnection(conn)
			}

		case subUnsub := <-this.subscribeCh:
			subjects, conn := subUnsub.subjects, subUnsub.Conn
			for _, subject := range subjects {
				if _, ok := this.subjectSubscribers[subject]; !ok {
					this.subjectSubscribers[subject] = make(map[*Connection]bool)
				}
				this.subjectSubscribers[subject][conn] = true
				this.connectionSubjects[conn.Id][subject] = true
			}

		case subUnsub := <-this.unsubscribeCh:
			subjects, conn := subUnsub.subjects, subUnsub.Conn
			for _, subject := range subjects {
				this.removeConnectionFromSubjectSubscribers(subject, conn)
				delete(this.connectionSubjects[conn.Id], subject)
			}

		case message := <-this.messageCh:
			if connections, ok := this.subjectSubscribers[message.Subject]; ok {
				data, _ := json.Marshal(message)
				if this.config.SupportZip {
					data, _ = GzipEncode(data)
				}
				for conn := range connections {
					go func(c *Connection) {
						logrus.Infof("conn%s c.SendCh <- data", c.Id)
						c.SendCh <- data
					}(conn)
				}
			}
		}
	}
}

func (this *ConnManager) Send(message Message) {
	this.messageCh <- &message
}

func (this *ConnManager) Subscribe(conn *Connection, subjects []string) {
	this.subscribeCh <- &SubUnSub{subjects: subjects, Conn: conn}
}

func (this *ConnManager) UnSubscribe(conn *Connection, subjects []string) {
	this.unsubscribeCh <- &SubUnSub{subjects: subjects, Conn: conn}
}

func (this *ConnManager) CreateConnection(Socket *websocket.Conn) *Connection {
	conn := &Connection{
		manager: this,
		Id:      uuid.New(),
		Socket:  Socket,
		SendCh:  make(chan interface{}, this.config.ConnectionSendBufSize),
	}

	this.registerCh <- conn
	return conn
}
