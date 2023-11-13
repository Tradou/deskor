package chat

import (
	"fmt"
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
	flagIndex := strings.Index(msg.Text, " ")
	if flagIndex == -1 {
		return msg.Text[len(prefixCommand):]
	}

	return msg.Text[len(prefixCommand):flagIndex]
}

func parseFlags(entryFlags string, cmdFlags []string) []Flag {
	var flags []Flag

	r := regexp.MustCompile(fmt.Sprintf(`%s(\w+)(?:=("[^"]*"|\S+))?`, regexp.QuoteMeta(prefixFlag)))
	matches := r.FindAllStringSubmatch(entryFlags, -1)

	for _, match := range matches {
		key := match[1]
		var value string

		if len(match) > 2 {
			value = strings.Trim(match[2], `"`)
		}

		for _, cmdFlag := range cmdFlags {
			if cmdFlag == key {
				flags = append(flags, Flag{Key: key, Value: value})
			}
		}
	}

	return flags
}

func ValidateMandatory(entryFlags []Flag, cmdFlags []string) bool {
	if len(cmdFlags) > len(entryFlags) {
		return false
	}
	cmdFlagsMap := make(map[string]struct{})
	for _, m := range cmdFlags {
		cmdFlagsMap[m] = struct{}{}
	}

	for _, flag := range entryFlags {
		if _, found := cmdFlagsMap[flag.Key]; !found {
			return false
		}
	}

	return true
}
