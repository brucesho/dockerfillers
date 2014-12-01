package main

import (
	"errors"
	"fmt"
	"github.com/brucesho/dockerfillers/utils"
	"path"
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

	imageIds, err := utils.GetImageIdsFromName(args[0])
	if err != nil {
		return err
	}

	if len(imageIds) == 0 {
		return fmt.Errorf("No matching image found: %s", args[0])
	}

	dockerInfo, err := utils.GetDockerInfo()
	if err != nil {
		return err
	}
	driverRootDir := dockerInfo.StorageDriver.RootDir

	switch dockerInfo.StorageDriver.Kind {
	case "aufs":
		for _, imageId := range imageIds {
			imageDiffDir := utils.AufsGetDiffDir(driverRootDir, imageId)
			parentDiffDirs, err := utils.AufsGetParentDiffDirs(driverRootDir, imageId)
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

	case "devicemapper":
		for _, imageId := range imageIds {

			parentImage, err := utils.GetImageParent(path.Dir(driverRootDir), imageId)
			if err != nil {
				return err
			}

			rootfsPath, containerId, err := utils.DeviceMapperGetRootFS(driverRootDir, imageId)
			if err != nil {
				return err
			}
			defer utils.DeviceMapperRemoveContainer(containerId)

			parentRootfsPath, parentContainerId, err := utils.DeviceMapperGetRootFS(driverRootDir, parentImage)
			if err != nil {
				return err
			}
			defer utils.DeviceMapperRemoveContainer(parentContainerId)

			changes, err := utils.ChangesDirs(rootfsPath, parentRootfsPath)
			if err != nil {
				return err
			}

			for _, change := range changes {
				fmt.Printf("%s\n", change.String())
			}
		}

	default:
		return fmt.Errorf("Error: storage driver %s is unsupported.\n", dockerInfo.StorageDriver.Kind)
	}

	return nil
}

func (cmdSet *CommandSet) CmdDiffsize(args ...string) error {
	if len(args) != 1 {
		fmt.Println("Usage: diffsize [IMAGE]")
		return nil
	}

	imageIds, err := utils.GetImageIdsFromName(args[0])
	if err != nil {
		return err
	}

	if len(imageIds) == 0 {
		return fmt.Errorf("No matching image found: %s", args[0])
	}

	dockerInfo, err := utils.GetDockerInfo()
	if err != nil {
		return err
	}
	driverRootDir := dockerInfo.StorageDriver.RootDir

	switch dockerInfo.StorageDriver.Kind {
	case "aufs":
		for _, imageId := range imageIds {
			imageDiffDir := utils.AufsGetDiffDir(driverRootDir, imageId)
			parentDiffDirs, err := utils.AufsGetParentDiffDirs(driverRootDir, imageId)
			if err != nil {
				return err
			}

			changes, err := utils.AufsGetChanges(parentDiffDirs, imageDiffDir)
			if err != nil {
				return err
			}

			var totalSize int64 = 0
			for _, change := range changes {
				fmt.Printf("%s\n", change.String())
				totalSize += change.Size
			}
			fmt.Printf("%d\n", totalSize)
		}

	case "devicemapper":
		for _, imageId := range imageIds {

			parentImage, err := utils.GetImageParent(path.Dir(driverRootDir), imageId)
			if err != nil {
				return err
			}

			rootfsPath, containerId, err := utils.DeviceMapperGetRootFS(driverRootDir, imageId)
			if err != nil {
				return err
			}
			defer utils.DeviceMapperRemoveContainer(containerId)

			parentRootfsPath, parentContainerId, err := utils.DeviceMapperGetRootFS(driverRootDir, parentImage)
			if err != nil {
				return err
			}
			defer utils.DeviceMapperRemoveContainer(parentContainerId)

			changes, err := utils.ChangesDirs(rootfsPath, parentRootfsPath)
			if err != nil {
				return err
			}

			var totalSize int64 = 0
			for _, change := range changes {
				fmt.Printf("%s\n", change.String())
				totalSize += change.Size
			}
			fmt.Printf("%d\n", totalSize)
		}

	default:
		return fmt.Errorf("Error: storage driver %s is unsupported.\n", dockerInfo.StorageDriver.Kind)
	}

	return nil
}
