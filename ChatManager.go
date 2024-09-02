package main

import (
	"fmt"
)

const (
	CHAT_LINES_MAX = int(1024)
)

type PrivateChat struct {
}

type ChatManager struct {
	channels map[string](*ChatChannel)
	privMsg  map[string](*PrivateChat)
	myNick   string
}

func NewChatManager() *ChatManager {
	return &ChatManager{}
}

func (cm *ChatManager) NewJoin(channelName string, userName string) {
	if userName == cm.myNick {
		_, exists := cm.channels[channelName]
		if exists {
			panic(fmt.Sprintf("error: user %s joined channel that they are already in", userName))
		}
		c := NewChatChannel(channelName)
		c.JoinUser(userName)
		cm.channels[channelName] = c
	} else {
		err := cm.channels[channelName].JoinUser(userName)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (cm *ChatManager) NewPart(channelName string, userName string, reason string) {
	err := cm.channels[channelName].PartUser(userName, reason)
	if err != nil {
		fmt.Print(err)
	}
}

func (cm *ChatManager) NewTopic(channelName string, topic string) {
	cm.channels[channelName].SetTopic(topic)
}
