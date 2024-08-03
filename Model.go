package main

type Model interface {
	SetChannel(channel string)

	GetChannel() string

	SetController(c ControllerInterface)

	ServerReplyParser(reply []byte)
}
