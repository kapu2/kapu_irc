package main

import (
	"fmt"
)

type PrivateChat struct {
	// map with only keys
	name string
	// ring buffer
	chatMessages          [CHAT_LINES_MAX]string
	bufToWritePtr         int
	filledBufferPositions int
}

func NewPrivateChat(name string) *PrivateChat {
	return &PrivateChat{name: name}
}

// TODO: below are same with ChatChannel, combine them somehow

func (cc *PrivateChat) AddPrivMsg(msg string, source string) {
	// TODO: broadcast starts with dollar, we are not handling server messages yet
	chatMsg := fmt.Sprintf("<%s>: %s", source, msg)
	cc.AddChannelMessage(chatMsg)
}

func (cc *PrivateChat) AddChannelMessage(msg string) {
	cc.chatMessages[cc.bufToWritePtr] = msg
	cc.bufToWritePtr++
	cc.bufToWritePtr %= CHAT_LINES_MAX
	if cc.filledBufferPositions < CHAT_LINES_MAX {
		cc.filledBufferPositions++
	}
}
