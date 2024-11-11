package main

import (
	"fmt"
)

type PrivateChat struct {
	name string
	// ring buffer
	chatMessages          [CHAT_LINES_MAX]string
	bufToWritePtr         int
	filledBufferPositions int
}

func NewPrivateChat(name string) *PrivateChat {
	return &PrivateChat{name: name}
}

func (cc *PrivateChat) ChangeNick(oldNick string, newNick string) {
	if oldNick == cc.name {
		cc.name = newNick
		msg := fmt.Sprintf("user: %s changed nick to %s", oldNick, newNick)
		cc.AddChannelMessage(msg)
	}
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

func (cc *PrivateChat) GetChatContent() string {
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

func (cc *PrivateChat) GetInfo() string {
	return "Private message with: " + cc.name
}

func (cc *PrivateChat) GetUsers() string {
	return cc.name
}

func (cc *PrivateChat) GetName() string {
	return cc.name
}
