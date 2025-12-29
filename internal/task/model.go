package task

import (
	"encoding/json"
	"os"
	"slices"
)

type Tasker interface {
	GetID() int
	SetID(int)
	GetTitle() string
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
func (t Task) GetID() int {
	return t.ID
}
func (t *Task) SetID(ID int) {
	t.ID = ID
}
func (t Task) GetTitle() string {
	return t.Title
}
func (t Task) IsCompleted() bool {
	return t.Completed
}
func (t *Task) Toggle() bool {
	t.Completed = !t.Completed
	return t.Completed
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
func (c CompositeTask) IsCompleted() bool {
	c.UpdateStatus()
	return c.Completed
}

//CompositeTask

type TaskList struct {
	Name  string   `json:"name"`
	Tasks []Tasker `json:"tasks"`
}

func (tl TaskList) Write(path string) error {
	data, err := json.MarshalIndent(tl, "", "\t")
	if err != nil {
		return nil
	}
	os.WriteFile(path, data, 0644)
	return nil
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
func (tl *TaskList) AddTask(id int, before bool) {
	tl.Tasks = append(
		tl.Tasks,
		&Task{ID: tasksCounter + 1, Title: "place holder", Completed: false},
	)
}
