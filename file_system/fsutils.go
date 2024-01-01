package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func GetDirs(path string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dirs := make([]fs.DirEntry, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs, nil
}

func CreateDir(basePath, session string) error {
	session = strings.TrimSpace(session)
	return os.Mkdir(fmt.Sprintf("%s/%s", basePath, session), 0755)
}
