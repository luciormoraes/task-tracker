package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
			if err := listTasks(nil); err != nil {
				exitErr(err)
			}
		} else if len(args) == 1 {
			filter, err := parseStatus(args[0])
			if err != nil {
				exitErr(err)
			}
			if err := listTasks(&filter); err != nil {
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
