package main

type ControllerInterface interface {
	SendCommand(msg []byte)
}
