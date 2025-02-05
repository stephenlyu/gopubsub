package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/stephenlyu/gopubsub/config"
	"github.com/stephenlyu/gopubsub/connmanager"
	"github.com/stephenlyu/gopubsub/message"
	"golang.org/x/net/websocket"
)

const DEFAULT_PORT = 7788
const DEFAULT_END_POINT = "/source"

type Server struct {
	connMgr  *connmanager.ConnManager
	port     int
	endPoint string
}

func NewServer(conf config.Config, port int, endPoint string) *Server {
	if port == 0 {
		port = DEFAULT_PORT
	}
	if endPoint == "" {
		endPoint = DEFAULT_END_POINT
	}

	return &Server{
		connMgr:  connmanager.NewConnManager(conf),
		port:     port,
		endPoint: endPoint,
	}
}

func (this *Server) handler(socket *websocket.Conn) {
	conn := this.connMgr.CreateConnection(socket)
	conn.Run()
	logrus.Infof("Server.handler conn %s disconnect.", conn.Id)
}

func (this *Server) Start() {
	go this.connMgr.Start()
	http.Handle(this.endPoint, websocket.Handler(this.handler))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", this.port), nil))
}

func (this *Server) Publish(msg message.Message) {
	this.connMgr.Send(msg)
}
