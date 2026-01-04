package main

import (
	"fmt"
	core "gondo/internal/Core"
	"gondo/internal/task"
)


func main() {
	var tl task.TaskList
	err := tl.Read("./tasksdd2.json")
	if err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	
	if err:= core.Init(&tl); err!= nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	tl.Write("./tasksdd2.json")
	
}


