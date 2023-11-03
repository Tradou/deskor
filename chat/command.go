package chat

import (
	"strings"
)

var PrefixCommand = "/"

func IsCommand(msg Message) bool {
	return strings.HasPrefix(msg.Text, PrefixCommand)
}

func isAnnouncement(msg Message) bool {
	return msg.Sender == "SERVER" && msg.SenderIp == ""
}

func ShouldBeEncrypt(msg Message) bool {
	return !(IsCommand(msg) || isAnnouncement(msg))
}

func ShouldBeDecrypt(msg Message) bool {
	return ShouldBeEncrypt(msg)
}

func Dispatch(msg Message) Message {
	switch msg.Text[1:] {
	case "ping":
		return ping(msg)
	default:
		return unknown(msg)
	}
}

func ping(msg Message) Message {
	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      "Pong",
		Connected: msg.Connected,
	}
}

func unknown(msg Message) Message {
	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      "This command does not exists",
		Connected: msg.Connected,
	}
}
