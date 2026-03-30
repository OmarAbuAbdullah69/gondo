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
	tl *task.TaskList
}

var (
	tasksViewStartLine    = 1
	tasksViewSelectedLine = 1
	selectedTask          task.Tasker

	unfoldCompTasks []int

	horizontalCap   int
	horizontalIndex = 0
)

type strTask struct {
	t *task.Tasker
	s string
}

func Init(tl *task.TaskList, global bool) error {
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
func Update(event string, msg string) {
	if tui.Width <= 0 {

		return
	}

	// next i should make styling happen here in this function

	h := headerStr()
	tui.SetHeader(h)

	switch event {
	case "up":
		if tasksViewSelectedLine > 1 {
			tasksViewSelectedLine--
		}
	case "down":
		if tasksViewSelectedLine < maxscroll {
			tasksViewSelectedLine++
		}
	case "right":
		if horizontalIndex < horizontalCap-tui.Width+16 {
			horizontalIndex++
		}
	case "left":
		if horizontalIndex > 0 {
			horizontalIndex--
		}
	case "f":
		if selectedTask != nil {
			if slices.Contains(unfoldCompTasks, selectedTask.GetID()) {
				index := slices.Index(unfoldCompTasks, selectedTask.GetID())
				unfoldCompTasks = slices.Delete(unfoldCompTasks, index, index+1)
			} else {
				unfoldCompTasks = append(unfoldCompTasks, selectedTask.GetID())
			}
		}
	case " ":
		if selectedTask != nil {
			core.tl.Toggle(selectedTask.GetID())
		}
	case "d":
		if selectedTask != nil {
			core.tl.Delete(selectedTask.GetID())
			if tasksViewSelectedLine > maxscroll-1 {
				tasksViewSelectedLine = maxscroll - 1
			}
		}
		task.UpdateTL(&core.tl.Tasks)
	case "a":
		if selectedTask != nil {
			core.tl.AddTask(selectedTask.GetID(), msg, false)
		}
	case "A":
		core.tl.AddTask(0, msg, true)
	case "p":
		if selectedTask != nil {
			core.tl.PushDown(selectedTask.GetID(), msg)
		}
	case "i":
		if selectedTask != nil {
			core.tl.InsertTask(selectedTask.GetID(), msg, false)
		}
	case "I":
		if selectedTask != nil {
			core.tl.InsertTask(selectedTask.GetID(), msg, true)
			tasksViewSelectedLine++
		}
	case "r":
		if selectedTask != nil {
			core.tl.RenameTask(selectedTask.GetID(), msg)
		}
	case "R":
		core.tl.Name = msg
		//
	}
	stl := taskToString(&(core.tl.Tasks), 0)
	maxscroll = len(stl)

	// styling the tasks string
	for i, t := range stl {

		if len(t.s) > horizontalIndex && len(t.s) < horizontalIndex+tui.Width-6 {
			t.s = t.s[horizontalIndex:]
		} else if len(t.s) >= horizontalIndex+tui.Width-6 {
			t.s = t.s[horizontalIndex : horizontalIndex+tui.Width-6]
		} else {
			t.s = ""
		}

		style := lipgloss.NewStyle()
		if (*t.t).IsCompleted() {
			style = style.Foreground(lipgloss.Color("#6cef2f"))
		}
		switch (*t.t).(type) {
		case *task.CompositeTask:
			style = style.Bold(true)
			if len(t.s) != 0 {
				t.s += lipgloss.NewStyle().Foreground(lipgloss.Color("#fc4103")).Render(" #")
			}
		}
		if i+1 == tasksViewSelectedLine {
			selectedTask = *t.t
			style = style.Background(lipgloss.Color("#6f5e73"))
		}
		stl[i].s = style.Render(t.s)
	}

	f := footerStr()
	tui.SetFooter(f)

	headerHeight := lipgloss.Height(h)
	footerHeight := lipgloss.Height(f)
	bodyHeight := tui.Height - (headerHeight + footerHeight) - 4
	tasksLinesCount := bodyHeight - 2

	v := tasksViewSelectedLine - tasksViewStartLine
	if v == tasksLinesCount+1 {
		tasksViewStartLine++
	} else if v < 0 {
		tasksViewStartLine += v
	}

	var sl []string
	for _, st := range stl {
		sl = append(sl, st.s)
	}

	if len(sl) == 0 {
		tui.SetStrList(sl)
	}
	if tasksViewStartLine+tasksLinesCount <= len(stl) {
		tui.SetStrList(sl[tasksViewStartLine-1 : tasksViewStartLine+tasksLinesCount])
	} else {
		if tasksViewStartLine == 0 {
			tui.SetStrList(sl)
		} else {
			tui.SetStrList(sl[tasksViewStartLine-1:])
		}
	}

}

func taskToString(tl *[]task.Tasker, indent int) []strTask {
	var stl []strTask
	horizontalCap = 0
	for i, t := range *tl {
		str := strings.Repeat("\t", indent) + strconv.Itoa(i+1) + ". "
		switch tt := t.(type) {
		case *task.Task:
			str += tt.Title
			stl = append(stl, strTask{s: str, t: &(*tl)[i]})
		case *task.CompositeTask:
			str += tt.Title
			stl = append(stl, strTask{s: str, t: &(*tl)[i]})
			if slices.Contains(unfoldCompTasks, tt.GetID()) {
				stl = append(stl, taskToString(&tt.SubTasks, indent+1)...)
			}
		}
		l := len(str)
		if l > horizontalCap {
			horizontalCap = l
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
	if selectedTask != nil {
		rHS := lipgloss.NewStyle()
		if selectedTask.IsCompleted() {
			rHS = rHS.Foreground(lipgloss.Color("#6cef2f"))
			rightHand += rHS.Render("Completed")
		} else {
			rHS = rHS.Foreground(lipgloss.Color("#fc4103"))
			rightHand += rHS.Render("Incompleted")
		}
		switch selectedTask.(type) {
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
			strconv.Itoa(tasksViewSelectedLine)+":"+strconv.Itoa(horizontalIndex),
		),
		lipgloss.PlaceHorizontal(footerWidth/3, lipgloss.Right, rightHand),
	)

	str = lipgloss.PlaceHorizontal(tui.Width-6, lipgloss.Center, str)

	return footerStyle.Render(str)
}
