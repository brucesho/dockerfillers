package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

func main() {
	maxDepthPtr := flag.Int("maxDepth", 10, "max depth of directory tree")
	maxFilesPerLayerPtr := flag.Int("maxFilesPerLayer", 10, "max number of files per layer in directory tree")
	maxFileSizePtr := flag.Int("maxFileSize", 1024, "max file size in bytes")
	dstDirPtr := flag.String("dir", ".", "destination directory")
	flag.String("help", "", "creates random files and directories in file system")

	flag.Parse()

	total, err := createFilesAndDirs(*dstDirPtr, *maxDepthPtr, *maxFilesPerLayerPtr, *maxFileSizePtr)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Total bytes written: %d\n", total)
	}
}

func createFilesAndDirs(dstDir string, maxDepth int, maxFilesPerLayer int, maxFileSize int) (int64, error) {
	rand.Seed(time.Now().Unix())
	return createFilesAndDirsRec(dstDir, maxDepth, maxFilesPerLayer, maxFileSize)
}

func createFilesAndDirsRec(dstDir string, maxDepth int, maxFilesPerLayer int, maxFileSize int) (int64, error) {
	var total int64 = 0

	if maxDepth > 0 {
		numOfFilesInLayer := rand.Intn(maxFilesPerLayer + 1)
		for i := 0; i < numOfFilesInLayer; i++ {
			createDir := rand.Intn(2) == 0
			if createDir {
				newDir := path.Join(dstDir, generateFilename(i))
				if err := os.Mkdir(newDir, os.ModePerm); err != nil {
					return total, err
				}
				if fileInfo, err := os.Stat(newDir); err != nil {
					return total, err
				} else {
					total += fileInfo.Size()
				}
				size, err := createFilesAndDirsRec(newDir, maxDepth-1, maxFilesPerLayer, maxFileSize)
				if err != nil {
					return total, err
				}
				total += size
			} else {
				fileSize := rand.Intn(maxFileSize)
				if err := createFile(path.Join(dstDir, generateFilename(i)), fileSize); err != nil {
					return total, err
				}
				total += int64(fileSize)
			}
		}
	}

	return total, nil
}

func generateFilename(i int) string {
	ab := "abcdefghijklmnopqrstuvwxyz"
	return fmt.Sprintf("%s_%d", string(ab[rand.Intn(len(ab))]), i)
}

func createFile(filepath string, size int) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// write to file in 1K blocks
	for bytesWritten := 0; bytesWritten < size; {
		bytesToWrite := 1024
		if bytesWritten+bytesToWrite > size {
			bytesToWrite = size - bytesWritten
		}
		if _, err := file.Write(generateRandomData(bytesToWrite)); err != nil {
			return err
		}
		bytesWritten += bytesToWrite
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func generateRandomData(size int) []byte {
	ab := "abcdefghijklmnopqrstuvwxyz"
	data := make([]byte, size)
	for i, _ := range data {
		data[i] = ab[rand.Intn(len(ab))]
	}

	return data
}
