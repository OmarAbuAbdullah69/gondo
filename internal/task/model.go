package task

import (
	"encoding/json"
	"os"
)

type Tasker interface {
	GetID() int
	GetTitle() string
	IsCompletee() bool
	Toggle() bool
}

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CompositeTask struct {
	Task
	SubTasks []Tasker `json:"subtasks"`
}

func (c *CompositeTask) UpdateStatus() {
	for _, st := range c.SubTasks {
		if !st.IsCompletee() {
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
func (t Task) GetTitle() string {
	return t.Title
}
func (t Task) IsCompletee() bool {
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
func (c CompositeTask) GetTitle() string {
	return c.Title
}
func (c CompositeTask) IsCompletee() bool {
	c.UpdateStatus()
	return c.Completed
}
func (c *CompositeTask) Toggle() bool {
	c.Completed = !c.Completed
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

func (tl *TaskList) Toggle(id int) bool{
	list, index, found := find(&tl.Tasks, id)
	if found {
		t := (*list)[index]
		t.Toggle()
	}
	return found
}
