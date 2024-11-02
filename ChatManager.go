package main

import (
	"fmt"
)

const (
	CHAT_LINES_MAX = int(1024)
)

type ChatManager struct {
	channels            map[string](*ChatChannel)
	privMsg             map[string](*PrivateChat)
	chatNumberToChannel []ChatWindow
	statusChat          *StatusChat
	openChatWindow      ChatWindow
	myNick              string
	observer            Observer
}

func NewChatManager() *ChatManager {
	mgr := ChatManager{}
	mgr.channels = make(map[string](*ChatChannel))
	mgr.privMsg = make(map[string](*PrivateChat))
	mgr.statusChat = NewStatusChat()
	mgr.openChatWindow = mgr.statusChat
	mgr.chatNumberToChannel = append(mgr.chatNumberToChannel, mgr.statusChat)

	return &mgr
}

func (cm *ChatManager) NewJoin(channelName string, userName string) {
	if userName == cm.myNick || cm.myNick == "" {
		// TODO: there probably is a smarter place to define your own nickname for the first time
		cm.myNick = userName

		_, exists := cm.channels[channelName]
		if exists {
			panic(fmt.Sprintf("error: user %s joined channel that they are already in", userName))
		}
		c := NewChatChannel(channelName)
		cm.chatNumberToChannel = append(cm.chatNumberToChannel, c)

		c.JoinUser(userName)
		cm.channels[channelName] = c
		if cm.openChatWindow == nil {
			cm.openChatWindow = c
		}
	} else {
		err := cm.channels[channelName].JoinUser(userName)
		if err != nil {
			fmt.Print(err)
		}
	}
	cm.NotifyIfChanged(channelName)

}

func (cm *ChatManager) NewPart(channelName string, userName string, reason string) {
	err := cm.channels[channelName].PartUser(userName, reason)
	if err != nil {
		fmt.Print(err)
	}
	cm.NotifyIfChanged(channelName)
}

func (cm *ChatManager) NewTopic(channelName string, topic string) {
	cm.channels[channelName].SetTopic(topic)
	cm.NotifyIfChanged(channelName)
}

func (cm *ChatManager) NewTopicInfo(channelName string, nick string, timestamp string) {
	cm.channels[channelName].SetTopicInfo(nick, timestamp)
	cm.NotifyIfChanged(channelName)
}

func (cm *ChatManager) NewNamesReply(symbol string, channelName string, names string) {
	cm.channels[channelName].NamesReply(symbol, names)
	cm.NotifyIfChanged(channelName)
}

func (cm *ChatManager) NewNamesReplyEnd(channelName string, endOfNames string) {
	cm.channels[channelName].NamesReplyEnd(endOfNames)
	cm.NotifyIfChanged(channelName)
}

func (cm *ChatManager) NewPrivMsg(targets []string, source string, msg string) {
	for _, target := range targets {
		channel, ok := cm.channels[target]
		if ok {
			channel.AddPrivMsg(msg, source)
			cm.NotifyIfChanged(channel.name)
		} else {
			var pc *PrivateChat
			pc, ok = cm.privMsg[target]
			if ok {
				pc.AddPrivMsg(msg, source)
			} else {
				pc = NewPrivateChat(source)
				cm.chatNumberToChannel = append(cm.chatNumberToChannel, pc)
				// TODO: some check whether its user or channel join/part bug
				pc.AddPrivMsg(msg, source)
			}
		}
	}
}
func (cm *ChatManager) RegisterObserver(observer Observer) {
	cm.observer = observer
}

func (cm *ChatManager) GetOpenChatWindow() string {
	if cm.openChatWindow != nil {
		return cm.openChatWindow.GetName()
	}
	return ""
	// error logging?
}

func (cm *ChatManager) NotifyIfChanged(channelName string) {
	if cm.openChatWindow != nil && cm.openChatWindow.GetName() == channelName {
		cm.observer.NotifyObserver("chat", cm.openChatWindow.GetChatContent())
		cm.observer.NotifyObserver("info", cm.openChatWindow.GetInfo())
		cm.observer.NotifyObserver("names", cm.openChatWindow.GetUsers())
	}
}

func (cm *ChatManager) NewStatusMessage(msg string) {
	cm.statusChat.AddChannelMessage(msg)
	cm.NotifyIfChanged(cm.statusChat.GetName())
}

func (cm *ChatManager) ChangeOpenChatWindow(nr int) {
	if nr >= 0 && nr < len(cm.chatNumberToChannel) {
		cm.openChatWindow = cm.chatNumberToChannel[nr]
		cm.NotifyIfChanged(cm.openChatWindow.GetName())
	}
}
