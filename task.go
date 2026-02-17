package main

import (
	"fmt"
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
