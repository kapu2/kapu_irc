package main

type View interface {
	GetConnectionInfo() (ipAndPort string, nick string)
	GetInput()
	SetController(ci ControllerInterface)
	StartView()
}
