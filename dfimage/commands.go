package main

import (
	"fmt"
)

func (cmdSet *CommandSet) CmdHelp(args ...string) {
	if len(args) == 0 || len(args) > 1 {
		printUsage(cmdSet, false)
	} else {
		command, exists := cmdSet.commands[args[0]]
		if !exists {
			fmt.Printf("No help on %s - command does not exist\n", args[0])
		} else {
			fmt.Printf("%s: %s\n", args[0], command.description)
		}
	}
}
