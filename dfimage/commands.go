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

func (cmdSet *CommandSet) CmdDiffchanges(args ...string) error {
	if len(args) != 1 {
		fmt.Println("Usage: diffchanges [IMAGE]")
		return nil
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
			return fmt.Errorf("No matching image found: %s", args[0])
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

	} else {
		return fmt.Errorf("Error: storage driver %s is unsupported.\n", dockerInfo.StorageDriver.Kind)
	}
	return nil
}

func (cmdSet *CommandSet) CmdDiffsize(args ...string) error {
	if len(args) != 1 {
		fmt.Println("Usage: diffsize [IMAGE]")
		return nil
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
			return fmt.Errorf("No matching image found: %s", args[0])
		}

		for _, imageId := range imageIds {
			imageDiffDir := utils.AufsGetDiffDir(aufsRootDir, imageId)
			parentDiffDirs, err := utils.AufsGetParentDiffDirs(aufsRootDir, imageId)
			if err != nil {
				return err
			}

			changesSize, err := utils.AufsGetChangesSize(parentDiffDirs, imageDiffDir)
			if err != nil {
				return err
			}

			fmt.Printf("%d\n", changesSize)
		}

	} else {
		return fmt.Errorf("Error: storage driver %s is unsupported.\n", dockerInfo.StorageDriver.Kind)
	}
	return nil
}
