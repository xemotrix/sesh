package creator

import (
	"io/fs"
	"regexp"
	filesystem "sesh/file_system"
	"sesh/tmux"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	globalStyle = lipgloss.NewStyle().
			AlignVertical(lipgloss.Bottom).
			Width(60).
			PaddingLeft(2)

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8C093")).
			PaddingLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C34043")).
			Bold(true)

	feedbackStyle = lipgloss.NewStyle().
			PaddingTop(2).
			PaddingLeft(2)

	validStyle = feedbackStyle.Copy().
			Foreground(lipgloss.Color("#98BB6C"))

	invalidStyle = feedbackStyle.Copy().
			Foreground(lipgloss.Color("#C34043"))
)

type model struct {
	input         textinput.Model
	base          string
	validity      sessionNameEval
	width         int
	height        int
	invalidNames  []string
	targetSession string
	err           error
}

type sessionNameEval int

const (
	seshNameEmpty sessionNameEval = iota
	seshNameValid
	seshNameRepeated
	seshNameInvalidChars
)

func (m model) validateSessionName(name string) sessionNameEval {
	name = strings.TrimSpace(name)
	if name == "" {
		return seshNameEmpty
	}
	regex := regexp.MustCompile(`^([A-Za-z])(\w|-)*$`)
	if slices.Contains(m.invalidNames, name) {
		return seshNameRepeated
	}
	if !regex.Match([]byte(name)) {
		return seshNameInvalidChars
	}
	return seshNameValid
}

func initialModel(base string, dirs []fs.DirEntry, sessions []string) model {

	invalidNames := []string{}
	invalidNames = append(invalidNames, sessions...)
	for _, d := range dirs {
		invalidNames = append(invalidNames, d.Name())
	}

	ti := textinput.New()
	ti.Placeholder = "your-session-name"
	ti.Focus()
	ti.Cursor.SetMode(cursor.CursorStatic)
	return model{
		input:        ti,
		base:         base,
		invalidNames: invalidNames,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		m.err = nil
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			value := m.input.Value()
			if m.validity != seshNameValid {
				break
			}
			if err := filesystem.CreateDir(m.base, value); err != nil {
				m.err = err
				return m, nil
			}
			if err := tmux.CreateSession(m.base, value); err != nil {
				m.err = err
				return m, nil
			}
			m.targetSession = value
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.validity = m.validateSessionName(m.input.Value())
	return m, cmd
}

func (m model) View() string {
	base := ""
	base += filterStyle.Render(m.input.View()) + "\n"
	switch m.validity {
	case seshNameEmpty:
		base += feedbackStyle.Render("\n")
	case seshNameValid:
		base += validStyle.Render("✔ Valid!")
	case seshNameRepeated:
		base += invalidStyle.Render("✖ Conflict with existing session or directory")
	case seshNameInvalidChars:
		base += invalidStyle.Render("✖ Session name is invalid")
	}

	lstr := globalStyle.Render(base)
	if m.err != nil {
		lstr += "\n" + errorStyle.Render("Error: "+m.err.Error())
	} else {
		lstr += "\n"
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lstr)
}
