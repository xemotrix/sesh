package killer

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xemotrix/sesh/internal/tmux"
)

const (
	MAXLEN = 35
)

var (
	globalStyle = lipgloss.NewStyle().
			AlignVertical(lipgloss.Top).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#C8C093")).
			Width(40).
			PaddingLeft(2)

	errorStyle = lipgloss.NewStyle().
			PaddingTop(1).
			Foreground(lipgloss.Color("#C34043")).
			Bold(true)

	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	currentSessionStyle = lipgloss.NewStyle().
				PaddingTop(2).
				PaddingLeft(4).
				Foreground(lipgloss.Color("#727169")).
				AlignVertical(lipgloss.Bottom)

	focusedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF9E3B"))
)

type model struct {
	width          int
	height         int
	current        string
	index          int
	help           help.Model
	selected       map[string]struct{}
	targetSessions []string
	sessions       []string
	err            error
}

func initialModel(sessions []string, current string) model {
	sessions = slices.DeleteFunc(sessions, func(s string) bool {
		return s == current
	})
	help := help.New()

	return model{
		sessions: sessions,
		current:  current,
		help:     help,
		selected: make(map[string]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) CursorUp() {
	if m.index > 0 {
		m.index--
	}
}
func (m *model) CursorDown() {
	if m.index < len(m.sessions)-1 {
		m.index++
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		globalStyle = lipgloss.NewStyle().
			MarginTop(m.height / 3).
			Inherit(globalStyle)
	case tea.KeyMsg:
		m.err = nil
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, keys.Down):
			m.CursorDown()
		case key.Matches(msg, keys.Up):
			m.CursorUp()
		case key.Matches(msg, keys.Select):
			if _, ok := m.selected[m.sessions[m.index]]; ok {
				delete(m.selected, m.sessions[m.index])
			} else {
				m.selected[m.sessions[m.index]] = struct{}{}
			}
		case key.Matches(msg, keys.Confirm):
			if len(m.selected) == 0 {
				return m, tea.Quit
			}
			sessions := make([]string, 0, len(m.selected))
			for s := range m.selected {
				sessions = append(sessions, s)
			}
			if err := tmux.KillSessions(sessions); err != nil {
				m.err = err
				return m, nil
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	base := "\n"
	for i, s := range m.sessions {
		selected := " "
		if _, ok := m.selected[s]; ok {
			selected = "✗"
		}
		s = fmt.Sprintf("[%s] %s", selected, s)

		if i == m.index {
			s = focusedItemStyle.Render("* " + s)
		} else {
			s = itemStyle.Render(s)
		}
		base += s + "\n"
	}

	base += currentSessionStyle.Render(fmt.Sprintf(" Current session: %s", m.current)) + "\n"

	lstr := globalStyle.Render(base)

	if m.err != nil {
		lstr += "\n" + errorStyle.Render("Error: "+m.err.Error())
	} else {
		lstr += "\n"
	}

	lstr += "\n" + m.help.View(keys)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, lstr)
}
