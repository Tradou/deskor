package chat

import (
	"net"
)

type Message struct {
	Sender    string `json:"sender"`
	SenderIp  string `json:"senderIp"`
	Text      string `json:"text"`
	Connected int    `json:"int"`
}

type Client struct {
	Conn     net.Conn
	Messages chan Message
	Messager
	Commander
}

type Disconnect struct {
	Client Client
}
