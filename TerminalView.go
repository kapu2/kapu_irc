package main

import (
	"bufio"
	"os"
)

type TerminalView struct {
	conIf ControllerInterface
}

func NewTerminalView() *TerminalView {
	return &TerminalView{}
}

func RemoveLast(s string) string {
	if len(s) > 0 {
		return s[:len(s)-1]
	} else {
		return s
	}
}
func (tv *TerminalView) SetController(ci ControllerInterface) {
	tv.conIf = ci
}

func (tv *TerminalView) GetConnectionInfo() (ipAndPort string, nick string) {
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
	nick, _ = r.ReadString('\t')
	// we expect eof error here
	if nick == "" {
		panic("Nick is empty")
	}
	return
}

func (*TerminalView) GetInput() string {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return str
}
