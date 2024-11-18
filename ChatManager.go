package main

import (
	"fmt"
)

const (
	CHAT_LINES_MAX = int(1024)
)

const (
	NICK_RANK_REGULAR = iota
	NICK_RANK_VOICE
	NICK_RANK_OPERATOR
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

	_, exists := cm.channels[channelName]

	if userName == cm.myNick && !exists {
		c := NewChatChannel(channelName)
		cm.chatNumberToChannel = append(cm.chatNumberToChannel, c)
		c.JoinUser(userName)
		cm.channels[channelName] = c
		// its a new channel, we change window to it
		cm.openChatWindowNumber = len(cm.chatNumberToChannel) - 1
	} else if exists {
		err := cm.channels[channelName].JoinUser(userName)
		if err != nil {
			fmt.Print(err)
		}
		if userName == cm.myNick {
			// change window to it, because it is our own join
			nr, err := cm.ChannelNameToWindowNumber(channelName)
			if err == nil {
				cm.openChatWindowNumber = nr
			} else {
				panic(err)
			}
		}
	} else {
		panic(fmt.Sprintf("panic: user %s joining channel %s (a channel we are not in, and user that is not us)", userName, channelName))
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
func (cm *ChatManager) NewQuit(userName string, reason string) {
	for _, channel := range cm.channels {
		channel.QuitUser(userName, reason)
	}
	cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
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

func (cm *ChatManager) NewWelcome(myNick string, message string) {
	cm.myNick = myNick
	cm.NewStatusMessage(message)
	cm.NotifyIfChanged(cm.statusChat.GetName())
}

func (cm *ChatManager) NewNick(oldNick string, newNick string) {
	if oldNick == cm.myNick {
		cm.myNick = newNick
	} else {
		// privmsg doesn't need to know if our own name is changed
		for _, privMsg := range cm.privMsg {
			privMsg.ChangeNick(oldNick, newNick)
		}
	}
	for _, channel := range cm.channels {
		channel.ChangeNick(oldNick, newNick)
	}
	cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
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
					if source == cm.myNick {
						_, myNickInChannel := channel.users[cm.myNick]
						if !myNickInChannel {
							cm.NewStatusMessage("you are not in channel " + channel.name)
						} else {
							channel.AddPrivMsg(msg, source)
						}
					} else {
						channel.AddPrivMsg(msg, source)
					}

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

func (cm *ChatManager) CloseOpenWindow() {
	if cm.GetOpenChatWindow().CanBeClosed(cm.myNick) {
		channelName := cm.chatNumberToChannel[cm.openChatWindowNumber].GetName()

		if cm.openChatWindowNumber < len(cm.chatNumberToChannel)-1 {
			cm.chatNumberToChannel = append(cm.chatNumberToChannel[0:cm.openChatWindowNumber],
				cm.chatNumberToChannel[cm.openChatWindowNumber+1:]...)
		} else {
			cm.chatNumberToChannel = cm.chatNumberToChannel[0:cm.openChatWindowNumber]
		}
		// only one of these deletes do anything
		delete(cm.channels, channelName)
		delete(cm.privMsg, channelName)
		cm.openChatWindowNumber--
	} else {
		if cm.openChatWindowNumber == 0 {
			cm.NewStatusMessage("Statuswindow cannot be closed")
		} else {
			cm.NewStatusMessage("Part from channel before closing the window")
		}
	}
	cm.NotifyIfChanged(cm.GetOpenChatWindow().GetName())
}

func (cm *ChatManager) ChannelNameToWindowNumber(channelName string) (int, error) {
	for i, c := range cm.chatNumberToChannel {
		if channelName == c.GetName() {
			return i, nil
		}
	}
	return -1, fmt.Errorf("error: channelName %s not found", channelName)
}
