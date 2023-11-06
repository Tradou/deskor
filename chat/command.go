package chat

const prefixCommand = "/"

type Commander interface {
	Dispatch(msg Message) Message
}

type Cmd struct{}

func Dispatch(msg Message) Message {
	functions := map[string]struct {
		fn          func(message Message) Message
		description string
	}{
		"ping": {
			fn:          callPing,
			description: "Play ping-pong",
		},
	}

	if entry, found := functions[msg.Text[1:]]; found {
		return entry.fn(msg)
	} else {
		return callUnknown(msg)
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
		Text:      "This command does not exists",
		Connected: msg.Connected,
	}
}

func callPing(msg Message) Message {
	ms := Cmd{}
	return ms.ping(msg)
}

func callUnknown(msg Message) Message {
	ms := Cmd{}
	return ms.unknown(msg)
}
