package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const tasksFile = "tasks.json"

func loadTasks() ([]Task, error) {
	b, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // file will be created on save
		}
		return nil, err
	}
	if len(b) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, fmt.Errorf("invalid %s JSON: %w", tasksFile, err)
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, b, 0644)
}
