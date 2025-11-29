package task

import (
	"encoding/json"
)

func UnMarshalTaskList(data []byte, tl *TaskList) error {
	var raw struct {
		Name  string            `json:"name"`
		Tasks []json.RawMessage `json:"tasks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	name := raw.Name
	taskers, err := unMarshalTaskerList(raw.Tasks)
	if err != nil {
		return err
	}
	*tl = TaskList{Name: name, Tasks: taskers}
	return nil
}

func unMarshalTaskerList(trr []json.RawMessage) ([]Tasker, error) {
	var trl []Tasker
	for _, raw := range trr {
		t, err := detectTask(raw)
		if err != nil {
			return nil, err
		}
		trl = append(trl, t)
	}
	return trl, nil
}

func detectTask(raw json.RawMessage) (Tasker, error) {
	var compRaw struct {
		Task
		SubTasks []json.RawMessage `json:"subtasks"`
	}
	if err := json.Unmarshal(raw, &compRaw); err != nil {
		return nil, nil
	}
	if compRaw.SubTasks != nil {
		taskers, err := unMarshalTaskerList(compRaw.SubTasks)
		if err != nil {
			return nil, nil
		}
		return &CompositeTask{Task: compRaw.Task, SubTasks: taskers}, nil
	} else {
		return &compRaw.Task, nil
	}
}
