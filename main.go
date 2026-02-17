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
		desc := strings.Join(args, " ")
		id, err := createTask(desc)
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task added successfully (ID: %d)\n", id)
	case "update":
		if len(args) < 2 {
			exitErr(errors.New(`usage: task-cli update <id> "new description"`))
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
			exitErr(errors.New(`usage: task-cli mark-in-progress <id>`))
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
			exitErr(errors.New(`usage: task-cli mark-done <id>`))
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
	case "list":
		if len(args) == 0 {
			err := listTasks(nil)
			if err != nil {
				exitErr(err)
			}
		} else if len(args) == 1 {
			filter, err := parseStatus(args[0])
			if err != nil {
				exitErr(err)
			}
			err = listTasks(&filter)
			if err != nil {
				exitErr(err)
			}
		} else {
			exitErr(errors.New(`usage: task-cli list [todo|in-progress|done]`))
		}
	case "delete":
		if len(args) != 1 {
			exitErr(errors.New(`usage: task-cli delete <id>`))
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitErr(fmt.Errorf("invalid id %q: must be an integer", args[0]))
		}
		err = deleteTask(id)
		if err != nil {
			exitErr(err)
		}
		fmt.Printf("Task deleted successfully (ID: %d)\n", id)
	default:
		printUsage()
		os.Exit(1)
	}
}

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

func markInProgress(id int) error {
	return setStatus(id, StatusInProgress)
}

func markDone(id int) error {
	return setStatus(id, StatusDone)
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

func parseStatus(s string) (Status, error) {
	switch s {
	case "todo":
		return StatusTodo, nil
	case "in-progress":
		return StatusInProgress, nil
	case "done":
		return StatusDone, nil
	default:
		return "", fmt.Errorf("invalid status %q: must be todo|in-progress|done", s)
	}
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
