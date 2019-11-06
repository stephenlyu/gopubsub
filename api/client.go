package api

import (
	"github.com/stephenlyu/gopubsub/message"
	"golang.org/x/net/websocket"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"github.com/stephenlyu/gopubsub/config"
	"io/ioutil"
	"compress/gzip"
	"bytes"
)

const HEARTBEAT_INTERVAL = time.Minute

type MessageCallback interface {
	OnMessage(message *message.Message, raw []byte)
	OnError(err error)
}

type Client struct {
	host string
	port int
	endPoint string
	secure bool

	callback MessageCallback

	ws *websocket.Conn
	SendCh chan *message.Message
	quitCh chan struct{}
}

func NewClient(host string, port int, endPoint string, callback MessageCallback) *Client {
	return &Client{
		host: host,
		port: port,
		endPoint: endPoint,
		callback: callback,
		SendCh: make(chan *message.Message),
		quitCh: make(chan struct{}),
	}
}

func (this *Client) SetSecure(secure bool) {
	this.secure = secure
}

func (this *Client) Start() error {
	var url string
	var origin string

	if this.secure {
		url = fmt.Sprintf("wss://%s:%d%s", this.host, this.port, this.endPoint)
		origin = fmt.Sprintf("https://%s:%d", this.host, this.port)
	} else {
		url = fmt.Sprintf("ws://%s:%d%s", this.host, this.port, this.endPoint)
		origin = fmt.Sprintf("http://%s:%d", this.host, this.port)
	}
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	this.ws = ws

	go this.Read()
	go this.Write()
	go this.heartbeat()

	return nil
}

func (this *Client) Stop() {
	close(this.quitCh)
	this.ws.Close()
}

func (this *Client) onError(err error) {
	if this.callback != nil {
		this.callback.OnError(err)
	}
}

func (this *Client) heartbeat() {
	for {
		select {
		case <- time.After(HEARTBEAT_INTERVAL):
			this.SendCh <- &message.Message{Subject: "PING", Data: time.Now().Unix()}
		case <- this.quitCh:
			break
		}
	}
}

func (this *Client) Read() {
	var raw []byte
	var message *message.Message
	for {
		err := websocket.Message.Receive(this.ws, &raw)
		if err != nil {
			this.onError(err)
			break
		}

		if config.DEFAULT_CONFIG.SupportZip {
			raw, err = GzipDecode(raw)
		}

		err = json.Unmarshal(raw, &message)
		if err != nil {
			this.onError(err)
			break
		}

		if this.callback != nil {
			this.callback.OnMessage(message, raw)
		}
	}
}

func (this *Client) Write() {
	for {
		select {
		case message, ok := <- this.SendCh:
			if !ok {
				return
			}

			err := websocket.JSON.Send(this.ws, message)
			if err != nil {
				this.onError(err)
				return
			}
		}
	}
}

func (this *Client) Subscribe(subjects []string) {
	this.SendCh <- &message.Message{Subject: "SUBSCRIBE", Data: strings.Join(subjects, ",")}
}

func (this *Client) UnSubscribe(subjects []string) {
	this.SendCh <- &message.Message{Subject: "UNSUBSCRIBE", Data: strings.Join(subjects, ",")}
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}