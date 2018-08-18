package api

import (
	"github.com/stephenlyu/gopubsub/connmanager"
	"net/http"
	"log"
	"fmt"
	"github.com/stephenlyu/gopubsub/message"
	"github.com/stephenlyu/gopubsub/config"
	"golang.org/x/net/websocket"
	"github.com/Sirupsen/logrus"
)

type Server struct {
	connMgr *connmanager.ConnManager
	port int
	endPoint string
}

func NewServer(conf config.Config, port int, endPoint string) *Server {
	return &Server{
		connMgr: connmanager.NewConnManager(conf),
		port: port,
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
