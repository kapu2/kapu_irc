package main

import (
	"fmt"
	"net"
	"strings"
)

type Controller struct {
	viewInterface  View
	modelInterface Model
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

func (c *Controller) ParseInput(buf string) (stuffToSend []byte) {

	if strings.Index(buf, "/j") == 0 {
		buf = strings.Replace(buf, "/j", "JOIN", 1)
	} else {
		buf = "PRIVMSG " + c.modelInterface.GetChannel() + " :" + buf
	}

	stuffToSend = append([]byte(buf), "\r\n"...)
	return
}
func (c *Controller) Commander(conn net.Conn) {
	for {
		buf := c.viewInterface.GetInput()
		stuffToSend := c.ParseInput(buf)
		conn.Write(stuffToSend)
		fmt.Println("kapu-irc: Sending: ", string(stuffToSend))
	}
}

func NewController(v View, m Model) *Controller {
	c := &Controller{viewInterface: v, modelInterface: m}
	c.viewInterface.SetController(c)
	c.modelInterface.SetController(c)
	return c
}
func (controller Controller) StartProgram() {

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
	go controller.Commander(conn)
	controller.Listener(conn)
}
