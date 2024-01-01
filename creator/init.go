package creator

import (
	"io/fs"
	"os"

	filesystem "github.com/xemotrix/sesh/file_system"
	"github.com/xemotrix/sesh/tmux"

	tea "github.com/charmbracelet/bubbletea"
)

func InitBubbleTea(path string) error {
	sessions, err := tmux.GetSessions()
	if err != nil {
		return err
	}

	dirs, err := filesystem.GetDirs(path)
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		initialModel(path, dirs, sessions),
		tea.WithAltScreen(),
	)
	m, err := p.Run()
	model := m.(model)
	return tmux.SwitchToSession(model.targetSession)
}

func getDirs(path string) ([]fs.DirEntry, error) {
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
