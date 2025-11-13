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
		SubTasks: []task.Tasker{t1, t2}}
	ct.UpdateStatus()

	t3 := task.Task{ID: 4, Title: "make an actual project", Completed: true}
	tl := task.TaskList{Name: "leaning programing",
		Tasks: []task.Tasker{ct, t3}}

	err := tl.Write("tasks.json")
	if err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	var tl2 task.TaskList
	err = tl2.Read("./tasks.json")
	if err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	tl2.Write("./tasks2.json")
}
