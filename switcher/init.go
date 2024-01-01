package switcher

import (
	fs "github.com/xemotrix/sesh/file_system"
	"github.com/xemotrix/sesh/tmux"

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

	p := tea.NewProgram(
		initialModel(path, dirs, sessions),
		tea.WithAltScreen(),
	)
	go func() {
		p.Send(tea.KeyMsg{Type: -1, Runes: []rune{'/'}, Alt: false})
	}()
	m, err := p.Run()
	model := m.(model)
	return tmux.SwitchToSession(model.targetSession)
}
