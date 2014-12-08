package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func GetImageIdsFromName(imageName string) ([]string, error) {
	nameTagPair := strings.SplitN(imageName, ":", 2)
	repoName := nameTagPair[0]
	var tag = ""
	if len(nameTagPair) > 1 {
		tag = nameTagPair[1]
	}

	imageIds := make([]string, 0)
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return imageIds, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if fields[0] == repoName && tag == "" || fields[0] == repoName && fields[1] == tag {
				imageIds = append(imageIds, fields[2])
			}
		}
	}

	return imageIds, nil
}

func GetImageParent(dockerRoot, imageId string) (string, error) {
	imageJsonBytes, err := ioutil.ReadFile(path.Join(dockerRoot, "graph", imageId, "json"))
	if err != nil {
		return "", err
	}

	var imageJson interface{}
	if err := json.Unmarshal(imageJsonBytes, &imageJson); err != nil {
		return "", err
	}

	m := imageJson.(map[string]interface{})
	parent, ok := m["parent"]
	if !ok {
		return "", fmt.Errorf("image %s has no parent", imageId)
	}

	return parent.(string), nil
}

func sameFsTime(a, b time.Time) bool {
	return a == b ||
		(a.Unix() == b.Unix() &&
			(a.Nanosecond() == 0 || b.Nanosecond() == 0))
}

func sameFsTimeSpec(a, b syscall.Timespec) bool {
	return a.Sec == b.Sec &&
		(a.Nsec == b.Nsec || a.Nsec == 0 || b.Nsec == 0)
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GetDirTreeSize(dir string) (int64, error) {
	var size int64 = 0

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error while walking the file system, path: %s, error:%s", path, err)
			return err
		}

		// Rebase path
		path, err = filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		path = filepath.Join("/", path)

		// Skip AUFS metadata
		if matched, err := filepath.Match("/.wh..wh.*", path); err != nil || matched {
			return err
		}

		file := filepath.Base(path)
		if strings.HasPrefix(file, ".wh.") {
			return nil
		}

		size += f.Size()
		return nil
	})

	return size, err
}

func GetChangesRelativeToParent(imageId string) ([]Change, error) {
	dockerInfo, err := GetDockerInfo()
	if err != nil {
		return nil, err
	}
	return getChangesRelativeToParentHelper(imageId, dockerInfo.StorageDriver.Kind, dockerInfo.StorageDriver.RootDir)
}

func getChangesRelativeToParentHelper(imageId, storageDriverKind, driverRootDir string) ([]Change, error) {
	switch storageDriverKind {
	case "aufs":
		imageDiffDir := AufsGetDiffDir(driverRootDir, imageId)
		parentDiffDirs, err := AufsGetParentDiffDirs(driverRootDir, imageId)
		if err != nil {
			return nil, err
		}

		return AufsGetChanges(parentDiffDirs, imageDiffDir)

	case "devicemapper":
		parentImage, err := GetImageParent(path.Dir(driverRootDir), imageId)
		if err != nil {
			return nil, err
		}

		rootfsPath, containerId, err := DeviceMapperGetRootFS(driverRootDir, imageId)
		if err != nil {
			return nil, err
		}
		defer DeviceMapperRemoveContainer(containerId)

		parentRootfsPath, parentContainerId, err := DeviceMapperGetRootFS(driverRootDir, parentImage)
		if err != nil {
			return nil, err
		}
		defer DeviceMapperRemoveContainer(parentContainerId)

		return ChangesDirs(rootfsPath, parentRootfsPath)

	default:
		return nil, fmt.Errorf("Error: storage driver %s is unsupported.\n", storageDriverKind)
	}

	return nil, nil
}
