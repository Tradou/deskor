package chat

import "fmt"

const prefixCommand = "/"

var fns = map[string]struct {
	fn          func(message Message) Message
	description string
}{
	"help": {
		fn:          callHelp,
		description: "Describe commands",
	},
	"ping": {
		fn:          callPing,
		description: "Play ping-pong",
	},
}

type Commander interface {
	Dispatch(msg Message) Message
}

type Cmd struct{}

func Dispatch(msg Message) Message {
	if entry, found := fns[msg.Text[len(prefixCommand):]]; found {
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

func callUnknown(msg Message) Message {
	ms := Cmd{}
	return ms.unknown(msg)
}
