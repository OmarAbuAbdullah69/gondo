package main

import (
	"fmt"
	"gondo/internal/task"
)

func main() {
	t1 := task.Task{ID: 1, Title: "learn baiscs", Completed: true}
	t2 := task.Task{ID: 2, Title: "make a simple project", Completed: true}

	ct := task.CompositeTask{
		Task:     task.Task{ID: 3, Title: "learn a language"},
		SubTasks: []task.Task{t1, t2}}
	ct.UpdateStatus()

	t3 := task.Task{ID: 4, Title: "make an actual project", Completed: true}
	tl := task.TaskList{Name: "leaning programing",
		Tasks: []task.Tasker{ct, t3}}
	
	for _, ts := range tl.Tasks {
		fmt.Printf("%s -> done:%v\n", ts.GetTitle(), ts)
	}
}
