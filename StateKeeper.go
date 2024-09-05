package main

import (
	"fmt"
)

type StateKeeper struct {
	channel string
	conIf   ControllerInterface
	cm      *ChatManager
}

func NewStateKeeper() *StateKeeper {
	cm := NewChatManager()
	return &StateKeeper{cm: cm}
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

func NumericReplyValidityCheck(msg *IRCMessage) error {
	var err error
	if msg.source.sourceName == "" {
		err = fmt.Errorf("error: numeric reply source is empty")
	}
	// "A numeric reply SHOULD contain the target of the reply as the first parameter of the message.
	// its not a MUST so we don't check
	return err
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
			if len(parsedReply.parameters) == 3 {
				// parameter[0] is client, we dont need it
				sk.cm.NewTopic(parsedReply.parameters[1], parsedReply.parameters[2])
			} else {
				print(fmt.Errorf("error: RPL_TOPIC reply, amount of parameters expected: 3 got %d", len(parsedReply.parameters)))
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
	}

}
