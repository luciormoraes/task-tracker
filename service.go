package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func createTask(desc string) (int, error) {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return 0, errors.New("description cannot be empty")
	}

	tasks, err := loadTasks()
	if err != nil {
		return 0, err
	}

	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	now := time.Now().UTC()
	task := Task{
		ID:          nextID,
		Description: desc,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, task)
	if err := saveTasks(tasks); err != nil {
		return 0, err
	}
	return nextID, nil
}

func updateTask(id int, desc string) error {
	if id <= 0 {
		return errors.New("id must be > 0")
	}
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return errors.New("description cannot be empty")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = desc
			tasks[i].UpdatedAt = now
			return saveTasks(tasks)
		}
	}
	return fmt.Errorf("task %d not found", id)
}

func deleteTask(id int) error {
	if id <= 0 {
		return errors.New("id must be > 0")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task %d not found", id)
}

func setStatus(id int, s Status) error {
	if id <= 0 {
		return errors.New("id must be > 0")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = s
			tasks[i].UpdatedAt = now
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task %d not found", id)
}

func listTasks(filter *Status) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	printed := false

	for _, t := range tasks {
		if filter == nil || t.Status == *filter {
			fmt.Printf("#%d [%s] %s\n", t.ID, t.Status, t.Description)
			printed = true
		}
	}

	if !printed {
		fmt.Println("No tasks found.")
	}

	return nil
}

func markInProgress(id int) error {
	return setStatus(id, StatusInProgress)
}

func markDone(id int) error {
	return setStatus(id, StatusDone)
}
