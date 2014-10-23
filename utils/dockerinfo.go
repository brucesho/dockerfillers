package utils

import (
	"os/exec"
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
			fields := strings.Fields(lines[i+1])
			if fields[0] == "Root" && fields[1] == "Dir:" {
				dockerInfo.StorageDriver.RootDir = fields[2]
			}
		}
	}

	return dockerInfo, nil
}
