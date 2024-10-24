package main

import (
	"fmt"
)

type StateKeeper struct {
	conIf ControllerInterface
	cm    *ChatManager
}

func NewStateKeeper() *StateKeeper {
	cm := NewChatManager()
	return &StateKeeper{cm: cm}
}

func (sk *StateKeeper) SetController(ci ControllerInterface) {
	sk.conIf = ci
}

func (sk *StateKeeper) SetChatObserver(obs Observer) {
	sk.cm.RegisterObserver(obs)
}

func (sk *StateKeeper) GetOpenChatWindow() string {
	return sk.cm.GetOpenChatWindow()
}

func NumericReplyValidityCheck(msg *IRCMessage) error {
	var err error
	if msg.source.sourceName == "" {
		err = fmt.Errorf("error: numeric reply source is empty")
	}
	// "A numeric reply SHOULD contain the target of the reply as the first parameter of the message.
	// its not a MUST so we don't check
	return err
}

func (sk *StateKeeper) ServerReplyParser(reply string) {
	parsedReply, err := ParseIRCMessage(reply)
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
	} else if parsedReply.command == "JOIN" {
		if len(parsedReply.parameters) == 1 && parsedReply.source.sourceName != "" {
			sk.cm.NewJoin(parsedReply.parameters[0], parsedReply.source.sourceName)
		} else {
			if len(parsedReply.parameters) != 1 {
				err = fmt.Errorf("error: JOIN has unexpected amount of parameters expected: 1 got: %d", len(parsedReply.parameters))
			} else {
				err = fmt.Errorf("error: sourceName is empty while JOINing")
			}
			print(err.Error())
		}
	} else if parsedReply.command == RPL_TOPIC {
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 2 || len(parsedReply.parameters) == 3 {
				// parameter[0] is client, we dont need it
				if len(parsedReply.parameters) == 2 {
					sk.cm.NewTopic(parsedReply.parameters[1], "")
				} else {
					sk.cm.NewTopic(parsedReply.parameters[1], parsedReply.parameters[2])
				}
			} else {
				print(fmt.Errorf("error: RPL_TOPIC reply, amount of parameters expected: 2 or 3 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	} else if parsedReply.command == RPL_TOPICWHOTIME {
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 4 {
				// parameter[0] is client, we dont need it
				sk.cm.NewTopicInfo(parsedReply.parameters[1], parsedReply.parameters[2], parsedReply.parameters[3])
			} else {
				print(fmt.Errorf("error: RPL_TOPICWHOTIME reply, amount of parameters expected: 4 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	} else if parsedReply.command == RPL_NAMREPLY {
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 4 {
				// parameter[0] is client, we dont need it
				sk.cm.NewNamesReply(parsedReply.parameters[1], parsedReply.parameters[2], parsedReply.parameters[3])
			} else {
				print(fmt.Errorf("error: RPL_NAMREPLY reply, amount of parameters expected: 4 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	} else if parsedReply.command == RPL_ENDOFNAMES {
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 3 {
				// parameter[0] is client, we dont need it
				sk.cm.NewNamesReplyEnd(parsedReply.parameters[1], parsedReply.parameters[2])
			} else {
				print(fmt.Errorf("error: RPL_NAMREPLY reply, amount of parameters expected: 3 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	}
}
