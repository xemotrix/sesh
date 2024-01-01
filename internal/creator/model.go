package creator

import (
	"fmt"
	"io/fs"
	"regexp"
	"slices"

	filesystem "github.com/xemotrix/sesh/internal/file_system"
	"github.com/xemotrix/sesh/internal/tmux"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	WIDTH  = 50
	HEIGHT = 10
	MAXLEN = 35
)

var (
	globalStyle = lipgloss.NewStyle().
			AlignVertical(lipgloss.Top).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#C8C093")).
			Width(WIDTH).
			Height(HEIGHT).
			PaddingLeft(2)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8C093")).
			PaddingLeft(2).
			PaddingTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C34043")).
			Bold(true)

	tipStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingLeft(2).
			Foreground(lipgloss.Color("#727169"))

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
	seshNameTooLong
)

func (m model) validateSessionName(name string) sessionNameEval {
	// name = strings.TrimSpace(name)
	if name == "" {
		return seshNameEmpty
	}
	if len(name) > MAXLEN {
		return seshNameTooLong
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
	base += inputStyle.Render(m.input.View()) + "\n"

	switch m.validity {
	case seshNameEmpty:
		base += feedbackStyle.Render("\n")
		// base += feedbackStyle.Render("\n")
	case seshNameValid:
		base += validStyle.Render("✔ Valid!") + "\n"
	case seshNameRepeated:
		base += invalidStyle.Render("✖ Conflict with existing session or directory\n")
	case seshNameInvalidChars:
		base += invalidStyle.Render("✖ Session name is invalid\n")
	case seshNameTooLong:
		base += invalidStyle.Render("✖ Session name too long\n")
	default:
		panic("invalid session name evaluation")
	}

	lstr := globalStyle.Render(base) + "\n"

	if m.validity == seshNameValid {
		lstr += tipStyle.Render(fmt.Sprintf(" %s/%s", m.base, m.input.Value()))
	} else {
		lstr += "\n"
	}
	if m.err != nil {
		lstr += "\n" + errorStyle.Render("Error: "+m.err.Error())
	} else {
		lstr += "\n"
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lstr)
}
