package main

type ChatWindow interface {
	GetChatContent() string
	GetInfo() string
	GetUsers() string
	GetName() string
	CanBeClosed(myNick string) bool
}
