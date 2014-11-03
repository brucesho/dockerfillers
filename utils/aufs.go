package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func AufsGetDiffDir(aufsRootDir, imageId string) string {
	return path.Join(aufsRootDir, "diff", imageId)
}

func AufsGetParentDiffDirs(aufsRootDir, imageId string) ([]string, error) {
	parentIds, err := AufsGetParentIds(aufsRootDir, imageId)
	if err != nil {
		return nil, err
	}
	if len(parentIds) == 0 {
		return nil, fmt.Errorf("Dir %s does not have any parent layers", imageId)
	}
	layers := make([]string, len(parentIds))

	for i, p := range parentIds {
		layers[i] = path.Join(aufsRootDir, "diff", p)
	}

	return layers, nil
}

func AufsGetParentIds(aufsRootDir, imageId string) ([]string, error) {
	f, err := os.Open(path.Join(aufsRootDir, "layers", imageId))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := []string{}
	s := bufio.NewScanner(f)

	for s.Scan() {
		if t := s.Text(); t != "" {
			out = append(out, s.Text())
		}
	}
	return out, s.Err()
}

func AufsGetChanges(layers []string, rw string) ([]Change, error) {
	var changes []Change

	err := filepath.Walk(rw, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error while walking the file system, path: %s, error:%s", path, err)
			return err
		}

		// Rebase path
		path, err = filepath.Rel(rw, path)
		if err != nil {
			return err
		}
		path = filepath.Join("/", path)

		// Skip root
		if path == "/" {
			return nil
		}

		// Skip AUFS metadata
		if matched, err := filepath.Match("/.wh..wh.*", path); err != nil || matched {
			return err
		}

		change := Change{
			Path: path,
			Size: 0,
		}

		// Find out what kind of modification happened
		file := filepath.Base(path)
		// If there is a whiteout, then the file was removed
		if strings.HasPrefix(file, ".wh.") {
			originalFile := file[len(".wh."):]
			change.Path = filepath.Join(filepath.Dir(path), originalFile)
			change.Kind = ChangeDelete
			// find the original file and get its size
			for _, layer := range layers {
				stat, err := os.Stat(filepath.Join(layer, filepath.Dir(path), originalFile))
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				if err == nil {
					// found file, or maybe it's a directory?
					if stat.IsDir() {
						// TODO: compute size of directory
						dirSize, err := aufsGetDirTreeSize(filepath.Join(layer, filepath.Dir(path), originalFile))
						if err != nil {
							return err
						}
						change.Size -= dirSize
					} else {
						change.Size = -stat.Size()
					}
					break
				}
			}
		} else {
			// Otherwise, the file was added
			change.Kind = ChangeAdd
			change.Size = f.Size()

			// ...Unless it already existed in a top layer, in which case, it's a modification
			for _, layer := range layers {
				stat, err := os.Stat(filepath.Join(layer, path))
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				if err == nil {
					// The file existed in the top layer, so that's a modification

					// However, if it's a directory, maybe it wasn't actually modified.
					// If you modify /foo/bar/baz, then /foo will be part of the changed files only because it's the parent of bar
					if stat.IsDir() && f.IsDir() {
						if f.Size() == stat.Size() && f.Mode() == stat.Mode() && sameFsTime(f.ModTime(), stat.ModTime()) {
							// Both directories are the same, don't record the change
							return nil
						}
					}
					change.Kind = ChangeModify
					change.Size = f.Size() - stat.Size()
					break
				}
			}
		}

		// Record change
		changes = append(changes, change)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return changes, nil
}

func aufsGetDirTreeSize(dir string) (int64, error) {
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
