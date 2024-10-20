package main

type View interface {
	//GetConnectionInfo() (ipAndPort string, nick string)
	//GetInput()
	Observer
	SetController(ci ControllerInterface)
	StartView()
}
