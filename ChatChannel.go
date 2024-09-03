package main

import (
	"fmt"
)

type ChatChannel struct {
	// map with only keys
	name  string
	users map[string]struct{}
	topic string
	// ring buffer
	chatMessages          [CHAT_LINES_MAX]string
	bufToWritePtr         int
	filledBufferPositions int
}

func NewChatChannel(name string) *ChatChannel {
	return &ChatChannel{name: name, bufToWritePtr: 0, filledBufferPositions: 0}
}

func (cc *ChatChannel) JoinUser(user string) error {
	_, exists := cc.users[user]
	if exists {
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
	// TODO: how to find out who set the topic?
	msg := fmt.Sprintf("channel %s new topic set by ??? : %s", cc.name, cc.topic)
	cc.AddChannelMessage(msg)
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
