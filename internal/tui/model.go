package tui

import (
	"fmt"
	"gondo/internal/task"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tl     *task.TaskList
	input  textinput.Model
	listtening bool
	width  int
	hieght int
}

var (
	tasksLinesCount       int
	tasksViewPortCap      int
	tasksViewStartLine    = 1
	tasksViewSelectedLine = 1
	selectedTaskID        int

	unfoldCompTasks []int
)

func Newmodle(tl *task.TaskList) model {
	ti := textinput.New()
	ti.Placeholder = "enter text here"
	ti.CharLimit = 60
	return model{tl: tl, input: ti, listtening: false}
}

// the buffer is counted from top to bottom so by moving up we are dicreasing the value
func moveLineUp() bool {
	if tasksViewSelectedLine == 1 {
		return false
	}
	tasksViewSelectedLine--
	updateView()
	return true
}
func moveLineDwon() bool {
	if tasksViewSelectedLine == tasksLinesCount {
		return false
	}
	tasksViewSelectedLine++
	updateView()
	return true
}

// i seawr to God those thing confuse me
func updateView() {
	offset := tasksViewStartLine - tasksViewSelectedLine
	if offset == 1 {
		tasksViewStartLine--
	}
	if offset < -tasksViewPortCap {
		tasksViewStartLine++
	}
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
			moveLineUp()
		case "down":
			moveLineDwon()
		case "t":
			m.tl.Toggle(selectedTaskID)
		case "f":
				index,found := slices.BinarySearch(unfoldCompTasks, selectedTaskID)
			if found{
				unfoldCompTasks = slices.Delete(unfoldCompTasks, index, index+1)
			} else {
				unfoldCompTasks = append(unfoldCompTasks, selectedTaskID)
			}
		case "d":
			m.tl.Delete(selectedTaskID)
		case "A":
			m.tl.AddTask(selectedTaskID, false)
		case "a":
			m.tl.AddTask(selectedTaskID, true)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.hieght = msg.Height
	}
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

var taskLine = 0

func taskToString(tl []task.Tasker, indent int) []string {
	var ret []string
	for i, t := range tl {
		taskLine++
		if taskLine == tasksViewSelectedLine {
			selectedTaskID = t.GetID()
		}
		s := strings.Repeat("\t", indent) + fmt.Sprintf("%d. ", i+1) + t.GetTitle()
		switch tt := t.(type) {
		case *task.Task:
			if tt.IsCompleted() {
				ret = append(ret, taskCompletedStyle.Render(s))
			} else {
				ret = append(ret, taskIncompletedStyle.Render(s))
			}
		case *task.CompositeTask:
			if tt.IsCompleted() {
				ret = append(ret, ctaskCompletedStyle.Render(s))
			} else {
				ret = append(ret, ctaskIncompletedStyle.Render(s))
			}
			if slices.Contains(unfoldCompTasks, tt.GetID()) {
				ret = append(ret, taskToString(tt.SubTasks, indent+1)...)
			}
		}
	}
	return ret
}

func headerStr(m model) string {
	titleStyle := lipgloss.NewStyle()
	ret := lipgloss.Place(m.width-6, 0, lipgloss.Center, 0, titleStyle.Render(m.tl.Name))
	headerBorder := lipgloss.Border{Top: "",
		Bottom:      "󰇼",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: ""}
	headerStyle := lipgloss.NewStyle().
		Border(headerBorder).
		BorderBottomForeground(lipgloss.Color("#35fc03"))
	return headerStyle.Render(ret)
}

func footerStr(m model) string {
	ret := lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("model sel: %d", tasksViewSelectedLine),
		fmt.Sprintf("model view index: %d", tasksViewStartLine))
	ret = lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		ret,
		fmt.Sprintf("	 ID: %v  ", m.listtening),
	)
	footerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#fc4103"))
	return footerStyle.Render(ret)
}

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
	tasks := taskToString(m.tl.Tasks, 0)
	//
	footer := footerStr(m)
	header := headerStr(m)
	//culculate  body heightt
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	bodyHeight := m.hieght - (headerHeight + footerHeight) - 4

	// drawing the tasks
	taskLine = 0
	tasksLinesCount = len(tasks)
	tasksViewPortCap = bodyHeight - 1
	var tasksStr string
	lineStyle := lipgloss.NewStyle()
	for i := tasksViewStartLine; i <= len(tasks); i++ {
		// used for debuging
		/*if i == taskBufferStartLine && i == selectedLine {
			lineStyle = lineStyle.Background(lipgloss.Color("#ff9900"))
		} else if i == taskBufferStartLine {
			lineStyle = lineStyle.Background(lipgloss.Color("#ff0000"))
		} else */if i == tasksViewSelectedLine {
			lineStyle = lineStyle.Background(lipgloss.Color("#554e56"))
		}
		tasksStr += lineStyle.Render(tasks[i-1])
		lineStyle = lineStyle.Background(lipgloss.Color(""))
		if i-tasksViewStartLine != tasksViewPortCap {
			tasksStr += "\n"
			continue
		}
		break
	}
	body = lipgloss.Place(m.width-8, bodyHeight, lipgloss.Left, lipgloss.Top, tasksStr)
	view = lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinHorizontal(lipgloss.Left,lipgloss.Place(2, bodyHeight, lipgloss.Left, lipgloss.Top, " "), bodyStyle.Render(body)),
		footer,
	)

	return frame.Render(view)
}

func StartDisplay(tl *task.TaskList) error {
	p := tea.NewProgram(Newmodle(tl), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
