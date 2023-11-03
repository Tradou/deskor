package chat

import (
	"deskor/encrypt"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type Message struct {
	Sender    string `json:"sender"`
	SenderIp  string `json:"senderIp"`
	Text      string `json:"text"`
	Connected int    `json:"int"`
	IChat
}

type Client struct {
	Conn     net.Conn
	Messages chan Message
}

type Disconnect struct {
	Client Client
}

type IChat interface {
	DispatchMessage()
	EncodeMessage(sender, text string) ([]byte, error)
	SendMessage(conn io.Writer, message []byte) error
	ReceiveMessage(conn io.Reader) (string, error)
	DecodeMessage(message string) (Message, error)
}

func (c *Client) EncodeMessage(sender, text string) ([]byte, error) {
	encrypter := encrypt.EncrypterImpl{}

	cypherText, err := encrypter.Encrypt(text)
	if err != nil {
		fmt.Print("Error while encrypting message")
	}
	message := Message{
		Sender: sender,
		Text:   cypherText,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return append(messageJSON, '\n'), nil
}

func (c *Client) SendMessage(conn io.Writer, message []byte) error {
	_, err := conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ReceiveMessage(conn io.Reader) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	message := string(buffer[:n])
	return message, nil
}

func (c *Client) DecodeMessage(message string) (Message, error) {
	var receivedMessage Message
	encrypter := encrypt.EncrypterImpl{}
	var err error

	if err := json.Unmarshal([]byte(message), &receivedMessage); err != nil {
		return Message{}, err
	}

	receivedMessage.Text, err = encrypter.Decrypt(receivedMessage.Text)
	if err != nil {
		fmt.Print("Error while decrypting message")
	}

	return receivedMessage, nil
}
