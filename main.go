package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const tasksFile = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "add":
		if len(args) < 1 {
			exitErr(errors.New(`missing description: task-cli add "Buy groceries"`))
		}
		id, err := createTask(args[0])
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task added successfully (ID: %d)\n", id)
	case "update":
		if len(args) < 2 {
			exitErr(errors.New(`missing description or id: task-cli update <id> "new description"`))
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitErr(fmt.Errorf("invalid id %q: must be an integer", args[0]))
		}

		desc := strings.Join(args[1:], " ")
		err = updateTask(id, desc)
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task updated successfully (ID: %d)\n", id)
	case "mark-in-progress":
		if len(args) != 1 {
			exitErr(errors.New(`missing id:: task-cli mark-in-progress <id>`))
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitErr(fmt.Errorf("invalid id %q: must be an integer", args[0]))
		}
		err = markInProgress(id)
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task marked in-progress successfully (ID: %d)\n", id)
	case "mark-done":
		if len(args) != 1 {
			exitErr(errors.New(`missing id: task-cli mark-done <id>`))
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitErr(fmt.Errorf("invalid id %q: must be an integer", args[0]))
		}
		err = markDone(id)
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task marked done successfully (ID: %d)\n", id)
	// TODO: update <id> <description>
	// TODO: delete <id>
	// TODO: mark-in-progress <id>
	// TODO: mark-done <id>
	// TODO: list [todo|in-progress|done]
	default:
		printUsage()
		os.Exit(1)
	}
}

func createTask(desc string) (int, error) {
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
	found := false
	if id <= 0 {
		return errors.New("id must be > 0")
	}

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
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %d not found", id)
	}

	return saveTasks(tasks)
}

func markInProgress(id int) error {
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
			tasks[i].Status = StatusInProgress
			tasks[i].UpdatedAt = now
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task %d not found", id)
}

func markDone(id int) error {
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
			tasks[i].Status = StatusDone
			tasks[i].UpdatedAt = now
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task %d not found", id)
}

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

func printUsage() {
	fmt.Println(`Usage:
  task-cli add "description"
  task-cli update <id> "description"
  task-cli delete <id>
  task-cli mark-in-progress <id>
  task-cli mark-done <id>
  task-cli list
  task-cli list todo
  task-cli list in-progress
  task-cli list done`)
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
