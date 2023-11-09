package chat

import (
	"regexp"
	"strings"
)

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

func getCommand(msg Message) string {
	return msg.Text[len("a"):strings.Index(msg.Text, " ")]
}

func getFlags(msg Message) []string {
	r := regexp.MustCompile(`--"([^"]*)"|--([^ ]*)`)
	matches := r.FindAllStringSubmatch(msg.Text, -1)

	var flags []string
	for _, match := range matches {
		if match[1] != "" {
			flags = append(flags, match[1])
		} else {
			flags = append(flags, match[2])
		}
	}
	return flags
}
