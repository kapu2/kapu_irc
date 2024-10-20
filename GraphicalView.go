package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	infoAreaHeight   = 5
	inputFieldHeight = 15
	namesFieldWidth  = 15
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type Keymap = struct {
	quit  key.Binding
	enter key.Binding
}

func NewTextarea(placeHolder string) textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = placeHolder
	t.ShowLineNumbers = false
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	//t.EndOfBufferCharacter = 'Q'
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.InsertNewline.SetEnabled(false)
	//t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	//t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type GraphicalView struct {
	width               int
	height              int
	keymap              Keymap
	chatArea            textarea.Model
	namesArea           textarea.Model
	inputArea           textarea.Model
	infoArea            textarea.Model
	controllerInterface ControllerInterface
}

func NewGraphicalView() *GraphicalView {
	gv := GraphicalView{
		keymap: Keymap{
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
			),
			enter: key.NewBinding(
				key.WithKeys("enter"),
			),
		},
	}
	gv.chatArea = NewTextarea("Chat area")
	gv.namesArea = NewTextarea("Names area")
	gv.inputArea = NewTextarea("Input area")
	gv.infoArea = NewTextarea("Info area")
	gv.inputArea.Focus()

	return &gv
}

func (gv *GraphicalView) Init() tea.Cmd {
	return textarea.Blink
}

func (gv *GraphicalView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	gv.controllerInterface.HandleReceivedMessages()

	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gv.keymap.quit):
			/*
				for i := range m.inputs {
					m.inputs[i].Blur()
				}
			*/
			return gv, tea.Quit
		case key.Matches(msg, gv.keymap.enter):
			str := gv.inputArea.Value()
			//print(str)
			gv.controllerInterface.HandleInput(str)
			gv.inputArea.Reset()
		}
	case tea.WindowSizeMsg:
		gv.height = msg.Height
		gv.width = msg.Width
	}

	gv.SizeTextAreas()

	// Update all textareas
	newModel, cmd := gv.chatArea.Update(msg)
	cmds = append(cmds, cmd)
	gv.chatArea = newModel

	newModel, cmd = gv.namesArea.Update(msg)
	cmds = append(cmds, cmd)
	gv.namesArea = newModel

	newModel, cmd = gv.inputArea.Update(msg)
	cmds = append(cmds, cmd)
	gv.inputArea = newModel

	newModel, cmd = gv.infoArea.Update(msg)
	cmds = append(cmds, cmd)
	gv.infoArea = newModel

	return gv, tea.Batch(cmds...)
}

func (gv *GraphicalView) SizeTextAreas() {
	smallWidth := (gv.width * 4) / 5
	smallHeight := (gv.height * 4) / 5

	gv.chatArea.SetWidth(smallWidth - namesFieldWidth)
	gv.chatArea.SetHeight(smallHeight - inputFieldHeight - infoAreaHeight)

	gv.namesArea.SetWidth(namesFieldWidth)
	gv.namesArea.SetHeight(smallHeight - inputFieldHeight - infoAreaHeight)

	gv.inputArea.SetWidth(smallWidth)
	gv.inputArea.SetHeight(inputFieldHeight)

	gv.infoArea.SetWidth(smallWidth)
	gv.infoArea.SetHeight(infoAreaHeight)
}

func (gv *GraphicalView) View() string {
	var views []string
	views = append(views, gv.chatArea.View())
	views = append(views, gv.namesArea.View())
	inputView := gv.inputArea.View()
	infoView := gv.infoArea.View()

	joinedViews := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	allViews := lipgloss.JoinVertical(lipgloss.Center, joinedViews, inputView)
	allViews = lipgloss.JoinVertical(lipgloss.Center, infoView, allViews)
	return allViews
}

func (gv *GraphicalView) StartView() {
	if _, err := tea.NewProgram(gv, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
func (gv *GraphicalView) SetController(ci ControllerInterface) {
	gv.controllerInterface = ci
}
func (gv *GraphicalView) NotifyObserver(field string, data string) {
	switch field {
	case "names":
		gv.namesArea.SetValue(data)
	case "info":
		gv.infoArea.SetValue(data)
	case "chat":
		gv.chatArea.SetValue(data)
	}
}
