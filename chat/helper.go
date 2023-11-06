package chat

import "strings"

type CommandHelper interface {
	IsCommand(msg Message) bool
	isAnnouncement(msg Message) bool
	ShouldBeEncrypt(msg Message) bool
	ShouldBeDecrypt(msg Message) bool
}

func IsCommand(msg Message) bool {
	return strings.HasPrefix(msg.Text, prefixCommand)
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
