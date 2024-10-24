package main

import (
	"bufio"
	"os"

	"github.com/rivo/tview"
)

type TerminalView struct {
	app   *tview.Application
	conIf ControllerInterface
}

func NewTerminalView(app *tview.Application) *TerminalView {
	return &TerminalView{app: app}
}

func RemoveLast(s string) string {
	if len(s) > 0 {
		return s[:len(s)-1]
	} else {
		return s
	}
}

func (tv *TerminalView) StartView() {
	go tv.GetInput()
}

func (tv *TerminalView) SetController(ci ControllerInterface) {
	tv.conIf = ci
}

func GetConnectionInfo() (ipAndPort string, nick string) {
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

func (tv *TerminalView) GetInput() {
	reader := bufio.NewReader(os.Stdin)
	var str string
	var err error
	for {
		str, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		// ignore enter smashing
		if str != "\n" {
			tv.conIf.HandleInput(string(str))
		}
	}
}
