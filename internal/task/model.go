package task

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CompositeTask struct {
	Task
	SubTasks []Task `json:"subtasks"`
}

func (c *CompositeTask) UpdateStatus() {
	for _, st := range c.SubTasks {
		if !st.Completed {
			c.Completed = false
			return
		}
	}
	c.Completed = true
}

type Tasker interface {
	GetID() int
	GetTitle() string
	IsCompletee() bool
}

//   Task
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
	Tasks []Tasker `json:"-"`
}
