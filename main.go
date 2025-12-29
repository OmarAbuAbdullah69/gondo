package main

import (
	"fmt"
	"gondo/internal/task"
	"gondo/internal/tui"
)


func main() {
	var tl task.TaskList
	err := tl.Read("./tasksdd2.json")
	if err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}

	if err:=tui.StartDisplay(&tl); err!= nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	tl.Write("./tasksdd2.json")
	
}


