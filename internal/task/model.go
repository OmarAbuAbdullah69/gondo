package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type Tasker interface {
	GetID() int
	SetID(int)
	GetTitle() string
	SetTitle(string)
	IsCompleted() bool
}

type Task struct {
	ID        int    `json:"-"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CompositeTask struct {
	Task
	SubTasks []Tasker `json:"subtasks"`
}

func (c *CompositeTask) UpdateStatus() {
	for _, st := range c.SubTasks {
		if !st.IsCompleted() {
			c.Completed = false
			return
		}
	}
	c.Completed = true
}

// Task
func (c Task) GetID() int {
	return c.ID
}
func (c *Task) SetID(ID int) {
	c.ID = ID
}
func (c Task) GetTitle() string {
	return c.Title
}
func (c *Task) SetTitle(s string) {
	c.Title = s
}
func (c Task) IsCompleted() bool {
	return c.Completed
}
func (c *Task) Toggle() bool {
	c.Completed = !c.Completed
	return c.Completed
}

// Task

// CompositeTask
func (c CompositeTask) GetID() int {
	return c.ID
}
func (c *CompositeTask) SetID(ID int) {
	c.ID = ID
}
func (c CompositeTask) GetTitle() string {
	return c.Title
}
func (c *CompositeTask) SetTitle(s string) {
	c.Title = s
}
func (c CompositeTask) IsCompleted() bool {
	c.UpdateStatus()
	return c.Completed
}

//CompositeTask

type TaskList struct {
	Name  string   `json:"name"`
	Tasks []Tasker `json:"tasks"`
}

// what does this functin baiscly do is checking if there are CompositeTasks withno SubTasks and turn them into normal tasks
func UpdateTL(tl *[]Tasker) {
	for i, t := range *tl {
		switch tt := t.(type) {
		case *CompositeTask:
			if len(tt.SubTasks) == 0 {
				(*tl)[i] = &Task{Title: tt.Title, Completed: false, ID: tt.ID}
			} else {
				UpdateTL(&tt.SubTasks)
			}
		}
	}
}
func (tl TaskList) Write(path string) error {
	data, err := json.MarshalIndent(tl, "", "\t")
	if err != nil {
		return nil
	}
	dir := filepath.Dir(path)
	if err1 := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err1
	}
	if err1 := os.WriteFile(path, data, 0644); err != nil {
		return err1
	}
	return err
}

func (tl *TaskList) Read(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = UnMarshalTaskList(data, tl)
	if err != nil {
		return err
	}
	return nil
}

func find(tl *[]Tasker, id int) (*[]Tasker, int, bool) {
	var listRet *[]Tasker
	var index int
	found := false
	for i, t := range *tl {
		if t.GetID() == id {
			listRet = tl
			index = i
			found = true
			break
		}

		//might cause a bug bc of the switch
		switch tt := t.(type) {
		case *CompositeTask:
			listRet, index, found = find(&tt.SubTasks, id)
			if found {
				return listRet, index, found
			}
		}
	}
	return listRet, index, found
}

func (tl *TaskList) Toggle(id int) bool {
	list, index, found := find(&tl.Tasks, id)
	if found {
		t := (*list)[index]
		switch tt := t.(type) {
		case *Task:
			tt.Toggle()
		default:
		}
	}
	return found
}
func (tl *TaskList) Delete(id int) bool {
	list, index, found := find(&tl.Tasks, id)
	if found {
		*list = slices.Delete(*list, index, index+1)
	}
	return found
}
func (tl *TaskList) AddTask(id int, title string, atbase bool) {
	t := &Task{Title: title, Completed: false, ID: tasksCounter}
	if atbase {
		tl.Tasks = append(tl.Tasks, t)
		tasksCounter++
		return
	}
	list, _, found := find(&tl.Tasks, id)
	if found {
		*list = append(*list, t)
		tasksCounter++
	}
}
func (tl *TaskList) PushDown(id int, title string) {
	list, index, found := find(&tl.Tasks, id)
	if found {
		switch t := (*list)[index].(type) {
		case *CompositeTask:
			t.SubTasks = append(t.SubTasks, &Task{Title: title, Completed: false, ID: tasksCounter})
		case *Task:
			nct := CompositeTask{Task: *t}
			nct.SubTasks = append(nct.SubTasks, &Task{Title: title, Completed: false, ID: tasksCounter})
			(*list)[index] = &nct
		}
		tasksCounter++
	}

}
func (tl *TaskList) InsertTask(id int, title string, before bool) {
	list, index, found := find(&tl.Tasks, id)
	if found {
		if !before {
			index++
		}
		*list = slices.Insert(
			*list,
			index,
			Tasker(&Task{Title: title, Completed: false, ID: tasksCounter}),
		)
		tasksCounter++
	}
}
func (tl *TaskList) RenameTask(id int, title string) {
	list, index, found := find(&tl.Tasks, id)
	if found {
		(*list)[index].SetTitle(title)
	}
}
