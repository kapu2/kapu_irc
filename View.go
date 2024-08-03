package main

type View interface {
	GetConnectionInfo() (ipAndPort string, nick string)
	GetInput() string
	SetController(ci ControllerInterface)
}

/*
import (
	"bufio"
	"os"
)

type View struct {
}

func NewView() *TerminalView {
	return &TerminalView{}
}

func RemoveLast(s string) string {
	if len(s) > 0 {
		return s[:len(s)-1]
	} else {
		return s
	}
}

func (*TerminalView) GetConnectionInfo() (ipAndPort string, nick string) {
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
*/
