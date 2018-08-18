package config


const (
	DEFAULT_CONNECTION_SEND_BUF_SIZE = 10
	DEFAULT_CONN_MANAGER_MESSAGE_BUF_SIZE = 1024
)

type Config struct {
	ConnectionSendBufSize 	int
	ConnManagerMessageBufSize int
}


var DEFAULT_CONFIG = &Config{
	ConnectionSendBufSize: DEFAULT_CONNECTION_SEND_BUF_SIZE,
	ConnManagerMessageBufSize: DEFAULT_CONN_MANAGER_MESSAGE_BUF_SIZE,
}
