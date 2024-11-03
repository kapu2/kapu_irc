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
	//openChatWindow      ChatWindow
	openChatWindowNumber int
	myNick               string
	observer             Observer
}

func NewChatManager() *ChatManager {
	mgr := ChatManager{}
	mgr.channels = make(map[string](*ChatChannel))
	mgr.privMsg = make(map[string](*PrivateChat))
	mgr.statusChat = NewStatusChat()
	//mgr.openChatWindow = mgr.statusChat
	mgr.openChatWindowNumber = 0
	mgr.chatNumberToChannel = append(mgr.chatNumberToChannel, mgr.statusChat)

	return &mgr
}

func RemoveExtraInfoFromTarget(s string) (string, error) {
	// we dont care about % or @ ( message only to half-ops or ops ), we show it if we can
	if len(s) > 0 && (s[0] == '%' || s[0] == '@') {
		s = RemoveFirstRuneFromString(s)
		return RemoveExtraInfoFromTarget(s)
	}
	if len(s) > 0 {
		return s, nil
	} else {
		err := "Error: empty target string in PRIVMSG (or constits only of % and @)"
		return s, fmt.Errorf("%s", err)
	}
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

		// its a new channel, we change window to it
		cm.openChatWindowNumber = len(cm.chatNumberToChannel) - 1
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
		var err error
		target, err = RemoveExtraInfoFromTarget(target)
		if err == nil {
			if target[0] == '$' {
				// broadcast
				str := "BROADCAST from " + source + ": " + msg
				cm.NewStatusMessage(str)
				cm.NotifyIfChanged(cm.statusChat.GetName())
			} else {
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
					cm.NotifyIfChanged(pc.name)
				}
			}

		} else {
			cm.NewStatusMessage(err.Error())
		}
	}
}
func (cm *ChatManager) RegisterObserver(observer Observer) {
	cm.observer = observer
}

func (cm *ChatManager) GetOpenChatWindow() ChatWindow {
	if cm.openChatWindowNumber < len(cm.chatNumberToChannel) && cm.openChatWindowNumber >= 0 {
		return cm.chatNumberToChannel[cm.openChatWindowNumber]
	} else {
		return nil
	}
	// error logging?
}

func (cm *ChatManager) NotifyIfChanged(channelName string) {
	if cm.GetOpenChatWindow() != nil && cm.GetOpenChatWindow().GetName() == channelName {
		cm.observer.NotifyObserver("chat", cm.GetOpenChatWindow().GetChatContent())
		cm.observer.NotifyObserver("info", cm.GetOpenChatWindow().GetInfo())
		cm.observer.NotifyObserver("names", cm.GetOpenChatWindow().GetUsers())
	}
}

func (cm *ChatManager) NewStatusMessage(msg string) {
	cm.statusChat.AddChannelMessage(msg)
	cm.NotifyIfChanged(cm.statusChat.GetName())
}

func (cm *ChatManager) ChangeOpenChatWindow(nr int) {
	if nr >= 0 && nr < len(cm.chatNumberToChannel) {
		cm.openChatWindowNumber = nr
		cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
	} else {
		cm.NewStatusMessage("Error: requested channel number does not exist")
	}
}

func (cm *ChatManager) ChangeToNextChatWindow() {
	channelAmount := len(cm.chatNumberToChannel)
	cm.openChatWindowNumber = (cm.openChatWindowNumber + 1) % channelAmount
	cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
}

func (cm *ChatManager) ChangeToPreviousChatWindow() {
	channelAmount := len(cm.chatNumberToChannel)
	if cm.openChatWindowNumber != 0 {
		cm.openChatWindowNumber = (cm.openChatWindowNumber - 1) % channelAmount
	} else {
		cm.openChatWindowNumber = channelAmount - 1
	}
	cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
}
