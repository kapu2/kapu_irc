package main

func main() {
	//app := tview.NewApplication()
	//terminalView := NewTerminalView(app)
	graphicalView := NewGraphicalView()

	stateKeeper := NewStateKeeper()
	controller := NewController(graphicalView, stateKeeper)
	controller.StartProgram()
}
