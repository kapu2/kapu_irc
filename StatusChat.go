package main

type StatusChat struct {
	// ring buffer
	chatMessages          [CHAT_LINES_MAX]string
	bufToWritePtr         int
	filledBufferPositions int
}

func NewStatusChat() *StatusChat {
	return &StatusChat{}
}

func (cc *StatusChat) AddChannelMessage(msg string) {
	cc.chatMessages[cc.bufToWritePtr] = msg
	cc.bufToWritePtr++
	cc.bufToWritePtr %= CHAT_LINES_MAX
	if cc.filledBufferPositions < CHAT_LINES_MAX {
		cc.filledBufferPositions++
	}
}

func (cc *StatusChat) GetChatContent() string {
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

func (cc *StatusChat) GetInfo() string {
	return "Status window"
}

func (cc *StatusChat) GetUsers() string {
	return ""
}

func (cc *StatusChat) GetName() string {
	return "?StatusWindow" //TODO: need to make these unique somehow
}