package main

import (
	"fmt"
	"flag"
	"math/rand"
	"time"
	"os"
	"path"
)

func main() {
	maxDepthPtr :=  flag.Int("maxDepth", 10, "max depth of directory tree")
	maxFilesPerLayerPtr := flag.Int("maxFilesPerLayer", 10, "max number of files per layer in directory tree")
	maxFileSizePtr := flag.Int("maxFileSize", 1024, "max file size in bytes")
	dstDirPtr := flag.String("dir", ".", "destination directory")
	flag.String("help", "", "creates random files and directories in file system")

	flag.Parse()

	if err := createFilesAndDirs(*dstDirPtr, *maxDepthPtr, *maxFilesPerLayerPtr, *maxFileSizePtr); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func createFilesAndDirs(dstDir string, maxDepth int, maxFilesPerLayer int, maxFileSize int) error {
	rand.Seed(time.Now().Unix())
	return createFilesAndDirsRec(dstDir, maxDepth, maxFilesPerLayer, maxFileSize)
}

func createFilesAndDirsRec(dstDir string, maxDepth int, maxFilesPerLayer int, maxFileSize int) error {
	if maxDepth > 0 {
		numOfFilesInLayer := rand.Intn(maxFilesPerLayer + 1)
		for i := 0; i < numOfFilesInLayer; i++ {
			createDir := rand.Intn(2) == 0
			if createDir {
				newDir := path.Join(dstDir, generateFilename(i))
				if err := os.Mkdir(newDir, os.ModePerm); err != nil {
					return err
				}
				createFilesAndDirsRec(newDir, maxDepth - 1, maxFilesPerLayer, maxFileSize)
			} else {
				file, err := os.Create(path.Join(dstDir, generateFilename(i)))
				if err != nil {
					return err
				}
				if _, err := file.Write(generateFiledata(rand.Intn(maxFileSize))); err != nil {
					return err
				}
				if err := file.Close(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func generateFilename(i int) string {
	ab := "abcdefghijklmnopqrstuvwxyz"
	return fmt.Sprintf("%s_%d", string(ab[rand.Intn(len(ab))]), i)
}

func generateFiledata(size int) []byte {
	ab := "abcdefghijklmnopqrstuvwxyz"
	data := make([]byte, size)
	for i, _ := range data {
		data[i] = ab[rand.Intn(len(ab))]
	}

	return data
}
