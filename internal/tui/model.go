package tui

import (
	"slices"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
}

var (
	Width  int
	Height int

	updateFunc func(string, string)
	header     string
	strList    []string
	footer     string

	input       textinput.Model
	listtening  bool = false
	inputString string
	keyListen   []string
	ListenEvent string
)

func SetKeyListen(s []string) {
	keyListen = append(keyListen, s...)
}

func SetUpdater(u func(string, string)) {
	updateFunc = u
}

func SetHeader(h string) {
	header = h
}
func SetStrList(sl []string) {
	strList = sl
}
func SetFooter(f string) {
	footer = f
}
func Newmodle() model {
	input = textinput.New()
	return model{}
}

func (m model) Init() tea.Cmd {
	input.Cursor.SetMode(cursor.CursorBlink)
	return textinput.Blink
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var event string
	var cmd tea.Cmd
	if listtening {
		input, cmd = input.Update(msg)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		event = msg.String()
		switch msg.String() {

		case "esc":
			if listtening {
				listtening = false
				input.Blur()
				input.Reset()
			} else {
				return m, tea.Quit
			}
		case "q":
			if !listtening {
				return m, tea.Quit
			}
		case "enter":
			listtening = false
			input.Blur()
			inputString = input.Value()
			input.Reset()
		default:
			event = msg.String()
			if !listtening {
				if slices.Contains(keyListen, event) {
					listtening = true
					input.Focus()
					ListenEvent = event
				}
			}
		}
	case tea.WindowSizeMsg:
		Width = msg.Width
		Height = msg.Height
	}
	if !listtening {
		if len(ListenEvent) != 0 && len(inputString) != 0 {
			updateFunc(ListenEvent, inputString)
			ListenEvent = ""
			inputString = ""
		} else {
			updateFunc(event, inputString)
		}
	} else {
		updateFunc("", "")
	}
	return m, cmd
}

func (m model) View() string {
	if listtening {
		footer = input.View()
		fs := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#fc4103")).
		Width(Width - 8)

		footer = fs.Render(footer)
	}
	// Send the UI for rendering
	var view string
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#fc4103")).
		Padding(0, 1, 0, 1).Margin(0)

	bodyBorder := lipgloss.Border{Top: "",
		Bottom:      "",
		Left:        "▎",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: ""}
	var body string

	for i, s := range strList {
		if i != len(strList)-1 {
			s += "\n"
		}
		body += s
	}
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	bodyHeight := Height - (headerHeight + footerHeight) - 2

	bodyStyle := lipgloss.NewStyle().Border(bodyBorder)

	body = lipgloss.Place(0, bodyHeight-3, lipgloss.Left, lipgloss.Top, body)
	// drawing the tasks
	view = lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.Place(0, bodyHeight, lipgloss.Left, lipgloss.Top, " "),
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
