package main

import (
	"bufio"
	"os"
)

func RemoveLast(s string) string {
	if len(s) > 0 {
		return s[:len(s)-1]
	} else {
		return s
	}
}

func GetConnectionInfo() (ipAndPort string, channel string, nick string) {
	fi, err := os.Open("connection_information.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	r := bufio.NewReader(fi)

	ipAndPort, err = r.ReadString('\t')
	ipAndPort = RemoveLast(ipAndPort)
	if err != nil {
		panic(err)
	}
	channel, err = r.ReadString('\t')
	channel = RemoveLast(channel)
	if err != nil {
		panic(err)
	}
	nick, _ = r.ReadString('\t')
	// we expect eof error here
	if nick == "" {
		panic("Nick is empty")
	}
	return
}
