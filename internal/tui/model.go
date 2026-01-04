package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	input      textinput.Model
	listtening bool
	width      int
	hieght     int
}

var (
	updateFunc func()
	header  string
	strList []string
	footer  string
)

func SetUpdater(u func()) {
	updateFunc = u
}

func SetHeader(h *string) {
	header = *h
}
func SetStrList(sl []string) {
	strList = sl
}
func SetFooter(f *string) {
	footer = *f
}
func Newmodle() model {
	ti := textinput.New()
	ti.Placeholder = "enter text here"
	ti.CharLimit = 60
	return model{input: ti, listtening: false}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
		case "down":
		case "t":
		case "f":
		case "d":
		case "A":
		case "a":
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.hieght = msg.Height
	}
	updateFunc()
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

var taskStyle = lipgloss.NewStyle().Margin(0, 0, 0, 0).TabWidth(4)

var taskCompletedStyle = lipgloss.NewStyle().
	Inherit(taskStyle).
	Foreground(lipgloss.Color("#14b83d"))
var taskIncompletedStyle = lipgloss.NewStyle().Inherit(taskStyle)

var ctaskCompletedStyle = lipgloss.NewStyle().
	Inherit(taskStyle).
	Foreground(lipgloss.Color("#14b83d")).
	Bold(true)
var ctaskIncompletedStyle = lipgloss.NewStyle().Inherit(taskStyle).Bold(true)

func (m model) View() string {
	// Send the UI for rendering
	var view string
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#fc4103")).
		Padding(0, 1, 0, 1)

	bodyBorder := lipgloss.Border{Top: "",
		Bottom:      "",
		Left:        "▎",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: ""}
	var body string
	bodyStyle := lipgloss.NewStyle().Border(bodyBorder)

	// tasks is spoused to be below this code but it will cause some issues rendring the data
	// tasks := core.TasksToString()
	//
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	bodyHeight := m.hieght - (headerHeight + footerHeight) - 4

	// drawing the tasks
	view = lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.Place(2, bodyHeight, lipgloss.Left, lipgloss.Top, " "),
			bodyStyle.Render(body),
		),
		footer,
	)

	return frame.Render(view)
}

func StartDisplay() error {
	p := tea.NewProgram(Newmodle(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
