package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Controller struct {
	theView  *View
	theModel *Model
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
		if bytes.Contains(buf, []byte("PING")) {
			if bytes.Contains(buf, []byte(":")) {
				start := 0
				end := 0
				for i, c := range buf {
					if c == ':' {
						start = i
					} else if start != 0 && c == '\r' {
						end = i
						break
					}
				}
				answer := append([]byte("PONG "), buf[start:end]...)
				answer = append(answer, []byte("\r\n")...)
				conn.Write(answer)
			} else {
				conn.Write([]byte("PONG\r\n"))
			}
			fmt.Println("kapu-irc: sent PONG")
		}

	}
}

func (c *Controller) ParseInput(buf string) (stuffToSend []byte) {

	if strings.Index(buf, "/j") == 0 {
		buf = strings.Replace(buf, "/j", "JOIN", 1)
	} else {
		// TODO: fix it
		c.theModel.GetChannel()
		buf = "PRIVMSG " + c.theModel.GetChannel() + " :" + buf
	}

	stuffToSend = append([]byte(buf), "\r\n"...)
	return
}
func (c *Controller) Commander(conn net.Conn) {
	for {
		buf := c.theView.GetInput()
		stuffToSend := c.ParseInput(buf)
		conn.Write(stuffToSend)
		fmt.Println("kapu-irc: Sending: ", string(stuffToSend))
	}
}

func NewController() *Controller {
	v := NewView()
	m := NewModel()
	return &Controller{theView: v, theModel: m}
}
func (controller Controller) StartProgram() {

	controller.theView.GetConnectionInfo()

	ipAndPort, channel, nick := controller.theView.GetConnectionInfo()

	controller.theModel.SetChannel(channel)

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
