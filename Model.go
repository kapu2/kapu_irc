package main

type Model interface {
	//SetChannel(channel string)
	// name of channel or user
	GetOpenChatWindow() string
	ChangeChatWindow(int)

	SetController(c ControllerInterface)

	SetChatObserver(obs Observer)

	ServerReplyParser(reply string)

	NewStatusMessage(msg string)
}
