package main

type Model interface {
	//SetChannel(channel string)
	// name of channel or user
	GetOpenChatWindow() string
	ChangeChatWindow(int)
	ChangeToNextChatWindow()
	ChangeToPreviousChatWindow()

	SetController(c ControllerInterface)

	SetChatObserver(obs Observer)

	ServerReplyParser(reply string)

	NewStatusMessage(msg string)

	GetMyNick() string
}
