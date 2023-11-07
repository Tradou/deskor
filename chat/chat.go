package chat

import (
	"encoding/json"
	"net"
)

const serverName = "SERVER"

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

type Server struct {
	Clients   map[Client]bool
	Join      chan Client
	Leave     chan Disconnect
	Messages  chan Message
	Connected int
}

func NewServer() *Server {
	return &Server{
		Clients:  make(map[Client]bool),
		Join:     make(chan Client),
		Leave:    make(chan Disconnect),
		Messages: make(chan Message),
	}
}

func (s *Server) JoinClient(client Client) {
	s.Join <- client
}

func (s *Server) LeaveClient(client Client) {
	s.Leave <- Disconnect{Client: client}
}

func (s *Server) AddMessage(message Message) {
	s.Messages <- message
}

func (s *Server) GetConnectedCount() int {
	return s.Connected
}

func (s *Server) IncrementConnectedCount() {
	s.Connected++
}

func (s *Server) DecrementConnectedCount() {
	s.Connected--
}

func (s *Server) SendWelcomeMessage(client Client) error {
	welcomeMessage := Message{
		Sender:    serverName,
		SenderIp:  "",
		Text:      "Someone has arrived",
		Connected: s.GetConnectedCount(),
	}
	welcomeMessageJSON, _ := json.Marshal(welcomeMessage)
	_, err := client.Conn.Write(welcomeMessageJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Clean(client Client) {
	delete(s.Clients, client)
	s.LeaveClient(client)
	client.Conn.Close()
}
