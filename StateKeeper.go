package main

import (
	"fmt"
)

type StateKeeper struct {
	channel string
	conIf   ControllerInterface
}

func NewStateKeeper() *StateKeeper {
	return &StateKeeper{}
}

func (sk *StateKeeper) SetController(ci ControllerInterface) {
	sk.conIf = ci
}

func (sk *StateKeeper) SetChannel(channel string) {
	sk.channel = channel
}

func (sk *StateKeeper) GetChannel() string {
	return sk.channel
}

func (sk *StateKeeper) ServerReplyParser(reply []byte) {
	parsedReply, err := ParseIRCMessage(string(string(reply)))
	if err != nil {
		print(err.Error())
	}
	if parsedReply.command == "PING" {
		if len(parsedReply.parameters) == 1 {
			msg := IRCMessage{command: "PONG"}
			msg.AddParameter(parsedReply.parameters[0])
			clientReply := IRCMessageToString(msg)
			sk.conIf.SendCommand([]byte(string(clientReply)))
		} else {
			err = fmt.Errorf("error: PING has unexpected amount of parameters expected: 1 got: %d", len(parsedReply.parameters))
			print(err.Error())
		}
	}
}
