package main

type ControllerInterface interface {
	SendCommand(msg []byte)
	HandleInput(input []rune)
}
