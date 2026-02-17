package main

import (
	"fmt"
	"os"
)

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
