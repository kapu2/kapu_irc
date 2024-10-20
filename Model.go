package main

type Model interface {
	//SetChannel(channel string)
	// name of channel or user
	GetOpenChatWindow() string

	SetController(c ControllerInterface)

	SetChatObserver(obs Observer)

	ServerReplyParser(reply string)
}
