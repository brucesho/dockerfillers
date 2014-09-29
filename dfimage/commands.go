package main

import (
	"fmt"
	"errors"
)

func (cmdSet *CommandSet) CmdHelp(args ...string) error {
	if len(args) == 0 || len(args) > 1 {
		printUsage(cmdSet, false)
	} else {
		command, exists := cmdSet.commands[args[0]]
		if !exists {
			return errors.New(fmt.Sprintf("No help on %s - command does not exist\n", args[0]))
		} else {
			fmt.Printf("%s: %s\n", args[0], command.description)
		}
	}

	return nil
}

func (cmdSet *CommandSet) CmdDiffsize(args ...string) error {
	if len(args) != 1 {
		return errors.New("IMAGE is a required arg")
	}

	return nil
}
