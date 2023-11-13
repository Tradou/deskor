package chat

import (
	"fmt"
)

const prefixCommand = "/"
const prefixFlag = "--"

var c *Cmd

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
	flags := parseFlags(msg.Text, fns["announce"].flags)

	if !ValidateMandatory(flags, fns["announce"].mandatory) {
		return Message{
			Sender:    "SERVER",
			SenderIp:  "",
			Text:      fmt.Sprintf("Announcement command have to be called with %v", fns["announce"].mandatory),
			Connected: msg.Connected,
		}
	}

	return Message{
		Sender:   "SERVER",
		SenderIp: "",
		// refactor this, it works only because there's one flag on this command.
		Text:      fmt.Sprintf("ANNOUNCEMENT: %s", flags[0].Value),
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
	return c.help(msg)
}

func callPing(msg Message) Message {
	return c.ping(msg)
}

func callAnnounce(msg Message) Message {
	return c.announce(msg)
}

func callUnknown(msg Message) Message {
	return c.unknown(msg)
}
