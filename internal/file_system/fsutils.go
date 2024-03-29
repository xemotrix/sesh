package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func CleanPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(user.HomeDir, path[1:])
	}
	return path, nil
}

func GetDirs(path string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dirs := make([]fs.DirEntry, 0, len(files))
	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			dirs = append(dirs, file)
		}
	}
	return dirs, nil
}

func CreateDir(basePath, session string) error {
	session = strings.TrimSpace(session)
	return os.Mkdir(fmt.Sprintf("%s/%s", basePath, session), 0755)
}
