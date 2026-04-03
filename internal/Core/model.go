package core

import (
	"gondo/internal/task"
	"gondo/internal/tui"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var core struct {
	tl                    *task.TaskList
	tasksViewStartLine    int
	tasksViewSelectedLine int
	selectedTask          task.Tasker

	unfoldCompTasks []int

	horizontalCap   int
	horizontalIndex int
}

type strTask struct {
	t task.Tasker
	s string
}

func Init(tl *task.TaskList, global bool) error {
	core.tasksViewSelectedLine = 1
	core.tasksViewStartLine = 1

	if len(tl.Name) == 0 && len(tl.Tasks) == 0 {
		if global {
			tl.Name = "Global todo list"
		} else {
			tl.Name = "New todo list"
		}
	}
	core.tl = tl
	task.UpdateTL(&core.tl.Tasks)
	tui.SetUpdater(Update)
	tui.SetKeyListen([]string{"a", "A", "p", "i", "I", "r", "R"})
	return tui.StartDisplay()
}

var maxscroll = 999 // its set to this number  just roto hold scrolling untill  the next key press just to make handle lenght later
func Update(event string, msg string) tui.ViewModel{
	var vm tui.ViewModel
	if tui.Width <= 0 {
		return vm
	}

	h := headerStr()
	vm.Header = h

	handleInput(event, msg)

	stl := taskToString(&(core.tl.Tasks), 0)
	maxscroll = len(stl)

	stl = styleSTL(stl)


	f := footerStr()
	vm.Footer = f
	headerHeight := lipgloss.Height(h)
	footerHeight := lipgloss.Height(f)
	bodyHeight := tui.Height - (headerHeight + footerHeight) - 4
	tasksLinesCount := bodyHeight - 2

	v := core.tasksViewSelectedLine - core.tasksViewStartLine
	if v == tasksLinesCount+1 {
		core.tasksViewStartLine++
	} else if v < 0 {
		core.tasksViewStartLine += v
	}

	sl := make([]string, len(stl))
	for i := range stl {
		sl[i] = stl[i].s
	}

	if len(sl) == 0 {
		vm.StrList = sl
	}
	if core.tasksViewStartLine+tasksLinesCount <= len(stl) {
		vm.StrList = sl[core.tasksViewStartLine-1 : core.tasksViewStartLine+tasksLinesCount]
	} else {
		if core.tasksViewStartLine == 0 {
			vm.StrList = sl
		} else {
			vm.StrList = sl[core.tasksViewStartLine-1:]
		}
	}
	return vm
}

func handleInput(event string, msg string) {
	switch event {
	case "up":
		if core.tasksViewSelectedLine > 1 {
			core.tasksViewSelectedLine--
		}
	case "down":
		if core.tasksViewSelectedLine < maxscroll {
			core.tasksViewSelectedLine++
		}
	case "right":
		if core.horizontalIndex < core.horizontalCap-tui.Width+16 {
			core.horizontalIndex++
		}
	case "left":
		if core.horizontalIndex > 0 {
			core.horizontalIndex--
		}
	case "f":
		if core.selectedTask != nil {
			if slices.Contains(core.unfoldCompTasks, core.selectedTask.GetID()) {
				index := slices.Index(core.unfoldCompTasks, core.selectedTask.GetID())
				core.unfoldCompTasks = slices.Delete(core.unfoldCompTasks, index, index+1)
			} else {
				core.unfoldCompTasks = append(core.unfoldCompTasks, core.selectedTask.GetID())
			}
		}
	case " ":
		if core.selectedTask != nil {
			core.tl.Toggle(core.selectedTask.GetID())
		}
	case "d":
		if core.selectedTask != nil {
			core.tl.Delete(core.selectedTask.GetID())
			if core.tasksViewSelectedLine > maxscroll-1 {
				if maxscroll > 1 {
					core.tasksViewSelectedLine = maxscroll - 1
				}
			}
		}
		task.UpdateTL(&core.tl.Tasks)
	case "a":
		if core.selectedTask != nil {
			core.tl.AddTask(core.selectedTask.GetID(), msg, false)
		}
	case "A":
		core.tl.AddTask(0, msg, true)
	case "p":
		if core.selectedTask != nil {
			core.tl.PushDown(core.selectedTask.GetID(), msg)
			core.unfoldCompTasks = append(core.unfoldCompTasks, core.selectedTask.GetID())
		}
	case "i":
		if core.selectedTask != nil {
			core.tl.InsertTask(core.selectedTask.GetID(), msg, false)
		}
	case "I":
		if core.selectedTask != nil {
			core.tl.InsertTask(core.selectedTask.GetID(), msg, true)
			core.tasksViewSelectedLine++
		}
	case "r":
		if core.selectedTask != nil {
			core.tl.RenameTask(core.selectedTask.GetID(), msg)
		}
	case "R":
		core.tl.Name = msg
		//
	}

}
func styleSTL(stl []strTask) []strTask{
	for i, t := range stl {

		if len(t.s) > core.horizontalIndex && len(t.s) < core.horizontalIndex+tui.Width-6 {
			t.s = t.s[core.horizontalIndex:]
		} else if len(t.s) >= core.horizontalIndex+tui.Width-6 {
			t.s = t.s[core.horizontalIndex : core.horizontalIndex+tui.Width-6]
		} else {
			t.s = ""
		}

		style := lipgloss.NewStyle()
		if (t.t).IsCompleted() {
			style = style.Foreground(lipgloss.Color("#6cef2f"))
		}
		switch (t.t).(type) {
		case *task.CompositeTask:
			style = style.Bold(true)
			if len(t.s) != 0 {
				t.s += lipgloss.NewStyle().Foreground(lipgloss.Color("#fc4103")).Render(" #")
			}
		}
		if i+1 == core.tasksViewSelectedLine {
			core.selectedTask = t.t
			style = style.Background(lipgloss.Color("#6f5e73"))
		}
		stl[i].s = style.Render(t.s)
	}
	return stl
}

func taskToString(tl *[]task.Tasker, indent int) []strTask {
	var stl []strTask
	core.horizontalCap = 0
	for i, t := range *tl {
		str := strings.Repeat("\t", indent) + strconv.Itoa(i+1) + ". "
		switch tt := t.(type) {
		case *task.Task:
			str += tt.Title
			stl = append(stl, strTask{s: str, t: (*tl)[i]})
		case *task.CompositeTask:
			str += tt.Title
			stl = append(stl, strTask{s: str, t: (*tl)[i]})
			if slices.Contains(core.unfoldCompTasks, tt.GetID()) {
				stl = append(stl, taskToString(&tt.SubTasks, indent+1)...)
			}
		}
		l := len(str)
		if l > core.horizontalCap {
			core.horizontalCap = l
		}
	}
	return stl
}

func headerStr() string {
	bodyBorder := lipgloss.Border{Top: "",
		Bottom:      "󰇼",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: ""}
	HeaderStyle := lipgloss.NewStyle().Border(bodyBorder).
		BorderForeground(lipgloss.Color("#6cef2f")).
		Margin(0, 0, 1, 0).Padding(1)

	str := lipgloss.Place(tui.Width-8, 0, lipgloss.Center, lipgloss.Center, core.tl.Name)

	return HeaderStyle.Render(str)
}

func footerStr() string {
	footerWidth := tui.Width - 6
	footerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#fc4103"))

	rightHand := "Status: "
	leftHand := "Type: "
	if core.selectedTask != nil {
		rHS := lipgloss.NewStyle()
		if core.selectedTask.IsCompleted() {
			rHS = rHS.Foreground(lipgloss.Color("#6cef2f"))
			rightHand += rHS.Render("Completed")
		} else {
			rHS = rHS.Foreground(lipgloss.Color("#fc4103"))
			rightHand += rHS.Render("Incompleted")
		}
		switch core.selectedTask.(type) {
		case *task.Task:
			leftHand += "Task"
		case *task.CompositeTask:
			leftHand += "CompositeTask"
		}
	}
	str := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.PlaceHorizontal(footerWidth/3, lipgloss.Left, leftHand),
		lipgloss.PlaceHorizontal(
			footerWidth/3,
			lipgloss.Center,
			strconv.Itoa(core.tasksViewSelectedLine)+":"+strconv.Itoa(core.horizontalIndex),
		),
		lipgloss.PlaceHorizontal(footerWidth/3, lipgloss.Right, rightHand),
	)

	str = lipgloss.PlaceHorizontal(tui.Width-6, lipgloss.Center, str)

	return footerStyle.Render(str)
}
