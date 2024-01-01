package killer

import (
	"github.com/xemotrix/sesh/internal/tmux"

	tea "github.com/charmbracelet/bubbletea"
)

func InitBubbleTea() error {
	sessions, err := tmux.GetSessions()
	if err != nil {
		return err
	}

	currentSession, err := tmux.GetCurrentSession()
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		initialModel(sessions, currentSession),
		tea.WithAltScreen(),
	)
	go func() {
		p.Send(tea.KeyMsg{Type: -1, Runes: []rune{'/'}, Alt: false})
	}()
	m, err := p.Run()
	model := m.(model)
	return tmux.KillSessions(model.targetSessions)
}
