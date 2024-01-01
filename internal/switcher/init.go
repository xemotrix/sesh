package switcher

import (
	"fmt"

	fs "github.com/xemotrix/sesh/internal/file_system"
	"github.com/xemotrix/sesh/internal/tmux"

	tea "github.com/charmbracelet/bubbletea"
)

func InitBubbleTea(path string) error {
	sessions, err := tmux.GetSessions()
	if err != nil {
		return err
	}
	dirs, err := fs.GetDirs(path)
	if err != nil {
		return err
	}

	currentSession, err := tmux.GetCurrentSession()
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		initialModel(path, dirs, sessions, currentSession),
		tea.WithAltScreen(),
	)
	m, err := p.Run()
	model := m.(model)

	targetSession := model.targetSession
	fmt.Println(targetSession)
	if targetSession == "" {
		return nil
	}
	return tmux.SwitchToSession(targetSession)
}
