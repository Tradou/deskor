package chat

import (
	"deskor/encrypt"
	"encoding/json"
)

type Messager interface {
	EncodeMessage(sender, text string) ([]byte, error)
	SendMessage(message []byte) error
	ReceiveMessage() (string, error)
	DecodeMessage(message string) (Message, error)
}

func (c *Client) EncodeMessage(sender, text string) ([]byte, error) {
	cypher := &encrypt.AesEncrypt{}

	message := Message{
		Sender: sender,
		Text:   text,
	}

	if ShouldBeEncrypt(message) {
		cypherText, err := cypher.Encrypt(text)
		if err != nil {
			return nil, err
		}
		message.Text = cypherText
	}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	return append(messageJSON, '\n'), nil
}

func (c *Client) SendMessage(message []byte) error {
	_, err := c.Conn.Write(message)
	return err
}

func (c *Client) ReceiveMessage() (string, error) {
	buffer := make([]byte, 1024)
	n, err := c.Conn.Read(buffer)
	if err != nil {
		return "", err
	}

	message := string(buffer[:n])
	return message, nil
}

func (c *Client) DecodeMessage(message string) (Message, error) {
	var receivedMessage Message
	cypher := encrypt.AesEncrypt{}
	var err error

	if err := json.Unmarshal([]byte(message), &receivedMessage); err != nil {
		return Message{}, err
	}

	if ShouldBeDecrypt(receivedMessage) {
		receivedMessage.Text, err = cypher.Decrypt(receivedMessage.Text)
		if err != nil {
			return Message{}, err
		}
	}

	return receivedMessage, nil
}
