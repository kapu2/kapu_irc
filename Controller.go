package main

import (
	"fmt"
	"net"
	"strings"
)

type Controller struct {
	viewInterface  View
	modelInterface Model
	messagesToSend chan []byte
}

func StartsWithReply(haystack []byte, command []byte) bool {
	for i, v := range command {
		if i < len(haystack) {
			if v != haystack[i] {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (c *Controller) Listener(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Println("Server reply:", string(buf[:n]))
		c.modelInterface.ServerReplyParser(buf)
	}
}

func (c *Controller) Commander(conn net.Conn) {
	for buf := range c.messagesToSend {
		conn.Write(buf)
		fmt.Println("kapu-irc: Sending: ", string(buf))
	}
}

func NewController(v View, m Model) *Controller {
	c := &Controller{viewInterface: v, modelInterface: m, messagesToSend: make(chan []byte)}
	c.viewInterface.SetController(c)
	c.modelInterface.SetController(c)
	return c
}
func (controller *Controller) StartProgram() {

	controller.viewInterface.GetConnectionInfo()

	ipAndPort, nick := controller.viewInterface.GetConnectionInfo()

	conn, err := net.Dial("tcp", ipAndPort)
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}
	conn.Write([]byte("CAP LS 302\r\n"))

	nickStr := "NICK " + nick + "\r\n"
	conn.Write([]byte(nickStr))
	conn.Write([]byte("USER d * 0 :What is this even\r\n"))

	controller.viewInterface.StartView()
	go controller.Commander(conn)
	controller.Listener(conn)
}

func (controller *Controller) SendCommand(msg []byte) {
	controller.messagesToSend <- msg
}

func (controller *Controller) HandleInput(input string) {
	if len(input) > 0 {
		if input[0] == '/' {
			controller.HandleInternalCommand(input)
		} else {
			controller.SendChatMessage(input)
		}
	}
}

func (controller *Controller) HandleInternalCommand(cmd string) {
	if strings.Index(string(cmd), "/j") == 0 {
		cmds := strings.Split(cmd, " ")
		if len(cmds) == 2 || len(cmds) == 3 {
			// channels and potential passwords delimited by ","
			// example: "/j #channel,#channel2 password1,password2"
			msg := IRCMessage{}
			msg.command = "JOIN"
			msg.parameters = cmds[1:]
			stringMsg := IRCMessageToString(msg)
			controller.messagesToSend <- []byte(stringMsg)
		}
	} else if strings.Index(string(cmd), "/n") == 0 {
		cmds := strings.Split(cmd, " ")
		if len(cmds) == 2 || len(cmds) == 3 {
			msg := IRCMessage{}
			msg.command = "NAMES"
			msg.parameters = cmds[1:]
			stringMsg := IRCMessageToString(msg)
			controller.messagesToSend <- []byte(stringMsg)
		}
	} else if strings.Index(string(cmd), "/t") == 0 {
		cmds := strings.Split(cmd, " ")
		if len(cmds) == 2 || len(cmds) == 3 {
			msg := IRCMessage{}
			msg.command = "TOPIC"
			msg.parameters = cmds[1:]
			stringMsg := IRCMessageToString(msg)
			controller.messagesToSend <- []byte(stringMsg)
		}
	}
}

func (controller *Controller) SendChatMessage(chatMsg string) {
	currentChannel := controller.modelInterface.GetChannel()
	msg := IRCMessage{}
	msg.command = "PRIVMSG"
	msg.AddParameter(currentChannel)
	msg.AddParameter(chatMsg)

	stringMsg := IRCMessageToString(msg)
	controller.messagesToSend <- []byte(stringMsg)
}
