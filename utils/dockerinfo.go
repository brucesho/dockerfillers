package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type DockerStorageDriver struct {
	Kind    string
	RootDir string
}

type DockerInfo struct {
	StorageDriver DockerStorageDriver
}

func GetDockerInfo() (DockerInfo, error) {
	var dockerStorageDriver = DockerStorageDriver{
		Kind:    "",
		RootDir: "",
	}

	var dockerInfo = DockerInfo{
		StorageDriver: dockerStorageDriver,
	}

	out, err := exec.Command("docker", "info").Output()
	if err != nil {
		return dockerInfo, err
	}

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "Storage Driver") {
			dockerInfo.StorageDriver.Kind = strings.Fields(line)[2]
			switch dockerInfo.StorageDriver.Kind {
			case "aufs":
				fields := strings.Fields(lines[i+1])
				if fields[0] == "Root" && fields[1] == "Dir:" {
					dockerInfo.StorageDriver.RootDir = fields[2]
				}

			case "devicemapper":
				fields := strings.Fields(lines[i+3])
				if fields[0] == "Data" && fields[1] == "file:" {
					dockerInfo.StorageDriver.RootDir = filepath.Dir(filepath.Dir(fields[2]))
				}

			default:
				return dockerInfo, fmt.Errorf("Storage driver %s is not supported", dockerInfo.StorageDriver.Kind)
			}

			break
		}
	}

	if dockerInfo.StorageDriver.RootDir == "" {
		return dockerInfo, errors.New("Failed to detect storage driver root directory")
	}

	return dockerInfo, nil
}
