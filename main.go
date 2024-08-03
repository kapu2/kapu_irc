package main

func main() {
	terminalView := NewTerminalView()
	stateKeeper := NewStateKeeper()
	controller := NewController(terminalView, stateKeeper)
	controller.StartProgram()
}
