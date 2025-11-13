package task

import (
	"encoding/json"
	"os"
)

type Tasker interface {
	GetID() int
	GetTitle() string
	IsCompletee() bool
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
func (c Task) GetID() int {
	return c.ID
}
func (c CompositeTask) GetTitle() string {
	return c.Title
}
func (c Task) IsCompletee() bool {
	return c.Completed
}

// Task

// CompositeTask
func (c CompositeTask) GetID() int {
	return c.ID
}
func (c Task) GetTitle() string {
	return c.Title
}
func (c CompositeTask) IsCompletee() bool {
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
