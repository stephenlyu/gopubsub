package message

type Message struct {
	Subject string			`json:"subject"`
	Data interface{}		`json:"data"`
}
