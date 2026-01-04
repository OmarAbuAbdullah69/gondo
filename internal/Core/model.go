package core

import (
	"gondo/internal/task"
	"gondo/internal/tui"
)

var core struct {
	tl *task.TaskList
}

var (
	tasksLinesCount       int
	tasksViewPortCap      int
	tasksViewStartLine    = 1
	tasksViewSelectedLine = 1
	selectedTaskID        int

	unfoldCompTasks []int
)

func Init(tl *task.TaskList) error {
	core.tl = tl
 	tui.SetUpdater(Update)
	return tui.StartDisplay()
}

func Update() {
	h := headerStr()
	sl := taskToString(core.tl, 0)
	f := FooterStr()
	tui.SetHeader(&h)
	tui.SetStrList(sl)
	tui.SetFooter(&f)
}

func taskToString(tl *task.TaskList, indent int) []string {
	var sl []string
	// for i, t := range tl.Tasks {
	//
// }
	return sl
}

func headerStr() string {
	var ret string
	return ret
}

func FooterStr() string {
	var ret string
	return ret
}
