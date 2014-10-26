package main

import (
	"errors"
	"fmt"
	"github.com/brucesho/dockerfillers/utils"
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

func (cmdSet *CommandSet) CmdDiffchanges(args ...string) error {
	if len(args) != 1 {
		return errors.New("IMAGE is a required arg")
	}

	dockerInfo, err := utils.GetDockerInfo()
	if err != nil {
		return err
	}

	if dockerInfo.StorageDriver.Kind == "aufs" {
		aufsRootDir := dockerInfo.StorageDriver.RootDir
		imageIds, err := utils.GetImageIdsFromName(args[0])
		if err != nil {
			return err
		}

		if len(imageIds) == 0 {
			fmt.Printf("No matching images found\n")
			return nil
		}

		for _, imageId := range imageIds {
			imageDiffDir := utils.AufsGetDiffDir(aufsRootDir, imageId)
			parentDiffDirs, err := utils.AufsGetParentDiffDirs(aufsRootDir, imageId)
			if err != nil {
				return err
			}

			changes, err := utils.AufsGetChanges(parentDiffDirs, imageDiffDir)
			if err != nil {
				return err
			}

			for _, change := range changes {
				fmt.Printf("%s\n", change.String())
			}
		}

	}
	return nil
}
