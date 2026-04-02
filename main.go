package main

import (
	"fmt"
	core "gondo/internal/Core"
	"gondo/internal/task"
	"os"
)

var helpMessage = `-g  to open the global todo list
-h for help

(arrows) to select a task
(A) to add a main task
(a) to add a neihboring task to the selected
(p) to push a task under the selected task
(i) to insert a task after the selected one
(I) to insert a task before the selected one
(space) to toggole the selected task
(f) to fold/unfold the selected task
(d) to delete the selected
(r) to rename the task`

func main() {
	globalList := false
	var path string
	if len(os.Args) == 1 {
		path = "./.gondo"
	} else if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-g":
			cache, err := os.UserCacheDir()
			if err != nil {
				fmt.Printf("error of type %T\n", err)
				fmt.Println("description:\n\t", err)
				return
			}
			globalList = true
			path = cache + "/gondo/global"
		case "-h":
			fmt.Println(helpMessage)
			return
		default:
			path = os.Args[1] + "/.gondo"
		}
	}
	var tl task.TaskList
	tl.Read(path)

	if err := core.Init(&tl, globalList); err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
	if len(tl.Tasks) + len(tl.Name) == 0 {
		return
	}
	err := tl.Write(path)
	if err != nil {
		fmt.Printf("error of type %T\n", err)
		fmt.Println("description:\n\t", err)
	}
}
