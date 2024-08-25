package main

import (
	"bytes"
	"fmt"
)

type StateKeeper struct {
	channel []rune
	conIf   ControllerInterface
}

func NewStateKeeper() *StateKeeper {
	return &StateKeeper{}
}

func (sk *StateKeeper) SetController(ci ControllerInterface) {
	sk.conIf = ci
}

func (sk *StateKeeper) SetChannel(channel []rune) {
	sk.channel = channel
}

func (sk *StateKeeper) GetChannel() []rune {
	return sk.channel
}

func (sk *StateKeeper) ServerReplyParser(reply []byte) {
	parsedReply, err := ParseIRCMessage([]rune(string(reply)))
	if err != nil {
		print(err.Error())
	}
	if SpliceIsSame(parsedReply.command, []rune("PING")) {
		if len(parsedReply.parameters) == 1 {
			msg := IRCMessage{command: []rune("PONG")}
			msg.AddParameter(parsedReply.parameters[0])
			clientReply := IRCMessageToString(msg)
			sk.conIf.SendCommand([]byte(string(clientReply)))
		} else {
			err = fmt.Errorf("error: PING has unexpected amount of parameters expected: 1 got: %d", len(parsedReply.parameters))
			print(err.Error())
		}
	}

	if bytes.Contains(reply, []byte("PING")) {
		if bytes.Contains(reply, []byte(":")) {
			start := 0
			end := 0
			for i, c := range reply {
				if c == ':' {
					start = i
				} else if start != 0 && c == '\r' {
					end = i
					break
				}
			}
			answer := append([]byte("PONG "), reply[start:end]...)
			answer = append(answer, []byte("\r\n")...)
			sk.conIf.SendCommand(answer)
		} else {
			sk.conIf.SendCommand([]byte("PONG\r\n"))
		}
		fmt.Println("kapu-irc: sent PONG")
		// TODO: below is all wrong
	} else if StartsWithReply(reply, []byte("JOIN")) {
		start := 0
		end := 0
		for i, c := range reply {
			if c == ':' {
				start = i + 1
			} else if start != 0 && c == '\r' {
				end = i
				break
			}
		}
		sk.SetChannel([]rune(string(reply[start:end])))
	}
}
