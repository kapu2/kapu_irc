package main

import (
	"fmt"
)

type ChatChannel struct {
	// map with only keys
	name              string
	users             map[string]struct{}
	topic             string
	topicSetBy        string
	topicSetTimestamp string // TODO: its unix timestamp at the moment
	// ring buffer
	chatMessages          [CHAT_LINES_MAX]string
	bufToWritePtr         int
	filledBufferPositions int
}

func NewChatChannel(name string) *ChatChannel {
	return &ChatChannel{name: name, users: make(map[string]struct{}), bufToWritePtr: 0, filledBufferPositions: 0}
}

func (cc *ChatChannel) JoinUser(user string) error {
	_, exists := cc.users[user]
	if !exists {
		cc.users[user] = struct{}{}
	} else {
		return fmt.Errorf("error: joining user: %s to channel: %s that already is in channel", user, cc.name)
	}
	msg := fmt.Sprintf("user: %s joined %s", user, cc.name)
	cc.AddChannelMessage(msg)

	return nil
}

func (cc *ChatChannel) PartUser(user string, reason string) error {
	_, exists := cc.users[user]
	if !exists {
		delete(cc.users, user)
	} else {
		return fmt.Errorf("error: parting user: %s from channel: %s that is not in channel", user, cc.name)
	}
	msg := fmt.Sprintf("user: %s left %s, reason: \"%s\"", user, cc.name, reason)
	cc.AddChannelMessage(msg)
	return nil
}

func (cc *ChatChannel) SetTopic(topic string) {
	cc.topic = topic
	msg := fmt.Sprintf("channel %s topic : %s", cc.name, cc.topic)
	cc.AddChannelMessage(msg)
}

func (cc *ChatChannel) SetTopicInfo(nick string, timestamp string) {
	cc.topicSetBy = nick
	cc.topicSetTimestamp = timestamp
	msg := fmt.Sprintf("channel %s topic set by %s at %ss after unix epoch", cc.name, nick, timestamp)
	cc.AddChannelMessage(msg)
}

func (cc *ChatChannel) NamesReply(symbol string, names string) {
	// TODO: do something with symbol or not, seems unnecessary information
	msg := fmt.Sprintf("channel %s names: %s", cc.name, names)
	cc.AddChannelMessage(msg)
}

func (cc *ChatChannel) NamesReplyEnd(endOfNames string) {
	cc.AddChannelMessage(endOfNames)
}

func (cc *ChatChannel) AddPrivMsg(msg string, source string) {
	// TODO: broadcast starts with dollar, we are not handling server messages yet
	chatMsg := fmt.Sprintf("<%s>: %s", source, msg)
	cc.AddChannelMessage(chatMsg)
}

func (cc *ChatChannel) AddChannelMessage(msg string) {
	cc.chatMessages[cc.bufToWritePtr] = msg
	cc.bufToWritePtr++
	cc.bufToWritePtr %= CHAT_LINES_MAX
	if cc.filledBufferPositions < CHAT_LINES_MAX {
		cc.filledBufferPositions++
	}
}

func (cc *ChatChannel) GetChatContent() string {
	var ret string
	// TODO: there is no chat scrolling yet, so we return 10 newest lines
	var start int
	if cc.filledBufferPositions > cc.bufToWritePtr {
		if cc.bufToWritePtr < 10 {
			start = CHAT_LINES_MAX - 1 - (10 - cc.bufToWritePtr)
		} else {
			start = cc.bufToWritePtr - 10
		}
	} else {
		if cc.bufToWritePtr < 10 {
			start = 0
		} else {
			start = cc.bufToWritePtr - 10
		}
	}
	for start != cc.bufToWritePtr {
		ret += cc.chatMessages[start] + "\n"
		start = (start + 1) % CHAT_LINES_MAX
	}
	return ret
}

func (cc *ChatChannel) GetUsers() string {
	var ret string
	for user := range cc.users {
		ret += user + "\n"
	}
	return ret
}

func (cc *ChatChannel) GetInfo() string {
	return "Channel name: " + cc.name
}

func (cc *ChatChannel) GetName() string {
	return cc.name
}
