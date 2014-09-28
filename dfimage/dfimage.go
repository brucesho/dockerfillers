package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	cmdSet := NewCommandSet()
	cmdSet.RegisterCommands()

	if len(os.Args) < 2 || !cmdSet.CommandExists(os.Args[1]) {
		printUsage(cmdSet, true)
	}

	err := cmdSet.InvokeCommand(os.Args[1], os.Args[2:])
	if err != nil {
		fmt.Printf("Error in invoking %s: %s\n", os.Args[1], err)
	}
}

func printUsage(cmdSet *CommandSet, exit bool) {
	fmt.Printf("Usage: %s command [arg...]\n\n", path.Base(os.Args[0]))

	fmt.Println("Commands:")
	for _, cmd := range cmdSet.GetCommands() {
		fmt.Printf("   %-10s   %s\n", cmd.name, cmd.description)
	}

	if exit {
		os.Exit(1)
	}
}
