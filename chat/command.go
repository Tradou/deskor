package chat

import (
	"strings"
)

var PrefixCommand = "/"

func isCommand(msg Message) bool {
	return strings.HasPrefix(msg.Text, PrefixCommand)
}

func isAnnouncement(msg Message) bool {
	return msg.Sender == "SERVER" && msg.SenderIp == ""
}

func ShouldBeEncrypt(msg Message) bool {
	return isCommand(msg) || isAnnouncement(msg)
}

func ShouldBeDecrypt(msg Message) bool {
	return isCommand(msg) || isAnnouncement(msg)
}
