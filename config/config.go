package config


const (
	DEFAULT_CONNECTION_SEND_BUF_SIZE = 100
	DEFAULT_CONN_MANAGER_MESSAGE_BUF_SIZE = 10240
)

type Config struct {
	ConnectionSendBufSize 	int
	ConnManagerMessageBufSize int
	SupportZip bool
}


var DEFAULT_CONFIG = &Config{
	ConnectionSendBufSize: DEFAULT_CONNECTION_SEND_BUF_SIZE,
	ConnManagerMessageBufSize: DEFAULT_CONN_MANAGER_MESSAGE_BUF_SIZE,
	SupportZip: false,
}
