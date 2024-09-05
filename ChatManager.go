package main

import (
	"fmt"
)

const (
	CHAT_LINES_MAX = int(1024)
)

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

func (cm *ChatManager) NewTopicInfo(channelName string, nick string, timestamp string) {
	cm.channels[channelName].SetTopicInfo(nick, timestamp)
}

func (cm *ChatManager) NewNamesReply(symbol string, channelName string, names string) {
	cm.channels[channelName].NamesReply(symbol, names)
}

func (cm *ChatManager) NewNamesReplyEnd(channelName string, endOfNames string) {
	cm.channels[channelName].NamesReplyEnd(endOfNames)
}

func (cm *ChatManager) NewPrivMsg(targets []string, source string, msg string) {
	for _, target := range targets {
		channel, ok := cm.channels[target]
		if ok {
			channel.AddPrivMsg(msg, source)
		} else {
			var pc *PrivateChat
			pc, ok = cm.privMsg[target]
			if ok {
				pc.AddPrivMsg(msg, source)
			} else {
				pc = NewPrivateChat(source)
				// TODO: some check whether its user or channel join/part bug
				pc.AddPrivMsg(msg, source)
			}
		}
	}
}
