package main

import (
	"bytes"
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
		sk.SetChannel(string(reply[start:end]))
	}
}
