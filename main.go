package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

func Listener(conn net.Conn) {
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

func ParseInput(buf string) (stuffToSend []byte) {

	if strings.Index(buf, "/j") == 0 {
		buf = strings.Replace(buf, "/j", "JOIN", 1)
	} else {
		// TODO: fix it
		buf = "PRIVMSG #test-channel :" + buf
	}

	stuffToSend = append([]byte(buf), "\r\n"...)
	return
}
func Commander(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		buf, _ := reader.ReadString('\n')
		stuffToSend := ParseInput(buf)
		conn.Write(stuffToSend)
		fmt.Println("kapu-irc: Sending: ", string(stuffToSend))
	}
}

func main() {
	ipAndPort, _, nick := GetConnectionInfo()
	conn, err := net.Dial("tcp", ipAndPort)
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}
	conn.Write([]byte("CAP LS 302\r\n"))

	nickStr := "NICK " + nick + "\r\n"
	conn.Write([]byte(nickStr))
	conn.Write([]byte("USER d * 0 :What is this even\r\n"))
	go Commander(conn)
	Listener(conn)
}
