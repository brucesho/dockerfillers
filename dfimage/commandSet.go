package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Command struct {
	name        string
	description string
}

type CommandSet struct {
	commands map[string]Command
}

func NewCommandSet() *CommandSet {
	return &CommandSet{make(map[string]Command, 0)}
}

func (cmdSet *CommandSet) RegisterCommands() {
	for _, c := range [][]string{
		{"diffsize", "size of an image (in bytes), relative to its parent"},
		{"diffchanges", "changes of an image relative to its parent"},
		{"help", "lists available commands"},
	} {
		cmd := &Command{c[0], c[1]}
		cmdSet.commands[c[0]] = *cmd
	}
}

func (cmdSet *CommandSet) GetCommands() map[string]Command {
	return cmdSet.commands
}

func (cmdSet *CommandSet) CommandExists(cmdName string) bool {
	_, exists := cmdSet.commands[cmdName]
	return exists
}

func (cmdSet *CommandSet) InvokeCommand(cmdName string, args []string) error {
	methodName := "Cmd" + strings.ToUpper(cmdName[:1]) + strings.ToLower(cmdName[1:])
	method := reflect.ValueOf(cmdSet).MethodByName(methodName)

	if !method.IsValid() {
		return errors.New(fmt.Sprintf("command %s is not implemented", cmdName))
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	results := method.Call(in)
	if len(results) == 0 {
		return nil
	} else {
		result := results[0].Interface()
		if result == nil {
			return nil
		} else {
			err, ok := result.(error)
			if ok {
				return err
			} else {
				return nil
			}
		}
	}
}
