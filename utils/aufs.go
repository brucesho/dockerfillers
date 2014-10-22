package utils

import (
	"path"
)

func AufsGetDiffDir(aufsRootDir, imageId string) string {
	return path.Join(aufsRootDir, "diff", imageId)
}

