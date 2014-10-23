package utils

import (
	"os/exec"
	"strings"
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
