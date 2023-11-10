package chat

import (
	"fmt"
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

func getFlags(msg string, cmd Commands) ([]Flag, error) {
	var flags []Flag

	args := strings.Fields(msg)

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if len(arg) > 2 && arg[0:2] == prefixFlag {
			var key, value string

			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				key = parts[0]
				value = parts[1]
			} else {
				key = arg[2:]
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], prefixFlag) {
					value = args[i+1]
					i++
				}
			}
			for _, s := range cmd.flags {
				if s == key {
					flags = append(flags, Flag{Key: key, Value: value})
				} else {
					return nil, fmt.Errorf("flag not exists")
				}
			}
		}
	}

	return flags, nil
}
