package chat

import (
	"fmt"
)

const prefixCommand = "/"
const prefixFlag = "--"

type Commands struct {
	fn          func(message Message) Message
	description string
	flags       []string
	mandatory   []string
}

type Flag struct {
	Key   string
	Value string
}

var fns = make(map[string]Commands)

func init() {
	fns["help"] = Commands{
		fn:          callHelp,
		description: "Describe commands",
	}

	fns["ping"] = Commands{
		fn:          callPing,
		description: "Play ping-pong",
	}

	fns["announce"] = Commands{
		fn:          callAnnounce,
		description: "Make an announce",
		flags:       []string{"text"},
		mandatory:   []string{"text"},
	}
}

type Commander interface {
	Dispatch(msg Message) Message
}

type Cmd struct{}

func Dispatch(msg Message) Message {
	command := getCommand(msg)
	if entry, found := fns[command]; found {
		return entry.fn(msg)
	}
	return callUnknown(msg)
}

func (c *Cmd) help(msg Message) Message {
	helpMessage := ""
	for k, v := range fns {
		helpMessage += fmt.Sprintf("%s: %s\n", k, v.description)
	}

	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      helpMessage,
		Connected: msg.Connected,
	}
}

func (c *Cmd) ping(msg Message) Message {
	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      "Pong",
		Connected: msg.Connected,
	}
}

func (c *Cmd) announce(msg Message) Message {
	flags, err := getFlags(msg.Text, fns["announce"])
	if err != nil {
		return Message{
			Sender:   "SERVER",
			SenderIp: "",
			// refactor this, it works only because there's one flag on this command.
			Text:      fmt.Sprintf("ANNOUNCEMENT: %s", flags[0].Value),
			Connected: msg.Connected,
		}
	}

	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      "Pong",
		Connected: msg.Connected,
	}
}

func (c *Cmd) unknown(msg Message) Message {
	return Message{
		Sender:    "SERVER",
		SenderIp:  "",
		Text:      fmt.Sprintf("This command does not exists, type %shelp for ...help", prefixCommand),
		Connected: msg.Connected,
	}
}

func callHelp(msg Message) Message {
	ms := Cmd{}
	return ms.help(msg)
}

func callPing(msg Message) Message {
	ms := Cmd{}
	return ms.ping(msg)
}

func callAnnounce(msg Message) Message {
	ms := Cmd{}
	return ms.announce(msg)
}

func callUnknown(msg Message) Message {
	ms := Cmd{}
	return ms.unknown(msg)
}
