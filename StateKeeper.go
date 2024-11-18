package main

import (
	"fmt"
	"strings"
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
	return sk.cm.GetOpenChatWindow().GetName()
}
func (sk *StateKeeper) CloseOpenWindow() {
	sk.cm.CloseOpenWindow()
}

func (sk *StateKeeper) ChangeChatWindow(nr int) {
	sk.cm.ChangeOpenChatWindow(nr)
}
func (sk *StateKeeper) ChangeToNextChatWindow() {
	sk.cm.ChangeToNextChatWindow()
}
func (sk *StateKeeper) ChangeToPreviousChatWindow() {
	sk.cm.ChangeToPreviousChatWindow()
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
	} else if parsedReply.command == "PRIVMSG" || parsedReply.command == "NOTICE" {
		// TODO: NOTICE will work similarly to how PRIVMSG works, it could be shown differently
		if len(parsedReply.parameters) > 1 {
			message := parsedReply.parameters[len(parsedReply.parameters)-1]
			for i := 0; i < len(parsedReply.parameters)-1; i++ {
				sk.cm.NewPrivMsg(parsedReply.parameters[:len(parsedReply.parameters)-1],
					parsedReply.source.sourceName,
					message)
			}
		}
	} else if parsedReply.command == "NICK" {
		oldNick := parsedReply.source.sourceName
		// there can be unnecessary garbage after "!" and we don't care about it
		newNick := strings.Split(parsedReply.parameters[0], "!")[0]
		sk.cm.NewNick(oldNick, newNick)

	} else if parsedReply.command == "PART" {
		user := parsedReply.source.sourceName
		channel := parsedReply.parameters[0]
		reason := "no reason"
		if len(parsedReply.parameters) == 2 {
			reason = parsedReply.parameters[1]
		}
		sk.cm.NewPart(channel, user, reason)
	} else if parsedReply.command == "QUIT" {
		user := parsedReply.source.sourceName
		reason := "no reason"
		if len(parsedReply.parameters) == 1 {
			reason = parsedReply.parameters[0]
		}
		sk.cm.NewQuit(user, reason)
	} else if parsedReply.command == "TOPIC" {
		topic := ""
		if len(parsedReply.parameters) == 2 {
			topic = parsedReply.parameters[1]
		}
		sk.cm.NewTopic(parsedReply.parameters[0], topic)
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
	} else if parsedReply.command == RPL_WELCOME {
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 2 {
				sk.cm.NewWelcome(parsedReply.parameters[0], parsedReply.parameters[1])
			} else {
				print(fmt.Errorf("error: RPL_NAMREPLY reply, amount of parameters expected: 3 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	} else if parsedReply.command == RPL_LUSERCLIENT ||
		parsedReply.command == RPL_LUSEROP ||
		parsedReply.command == RPL_LUSERUNKNOWN ||
		parsedReply.command == RPL_LUSERCHANNELS ||
		parsedReply.command == RPL_LUSERME ||
		parsedReply.command == RPL_ADMINME ||
		parsedReply.command == RPL_ADMINLOC2 ||
		parsedReply.command == RPL_ADMINEMAIL ||
		parsedReply.command == RPL_TRYAGAIN ||
		parsedReply.command == RPL_LOCALUSERS ||
		parsedReply.command == RPL_GLOBALUSERS ||
		parsedReply.command == RPL_MOTD ||
		parsedReply.command == RPL_MOTDSTART ||
		parsedReply.command == RPL_ENDOFMOTD {
		// a bunch of messages sent on server connection, we don't need to care what they are
		err = NumericReplyValidityCheck(&parsedReply)
		if err == nil {
			if len(parsedReply.parameters) == 2 {
				sk.cm.NewStatusMessage(parsedReply.parameters[1])
			} else {
				print(fmt.Errorf("error: RPL_NAMREPLY reply, amount of parameters expected: 3 got %d", len(parsedReply.parameters)))
			}
		} else {
			print(err.Error())
		}
	} else {
		// unhandled commands
		// most of them are decent error replies, so lets just show them on statuswindow
		str := "Unhandled response: " + reply
		sk.NewStatusMessage(str)
	}
}

func (sk *StateKeeper) NewStatusMessage(msg string) {
	sk.cm.NewStatusMessage(msg)
}

func (sk *StateKeeper) GetMyNick() string {
	return sk.cm.myNick
}
