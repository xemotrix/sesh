package switcher

import (
	"fmt"
	"io/fs"
	"slices"

	"github.com/sahilm/fuzzy"
	"github.com/xemotrix/sesh/internal/tmux"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MAX_ELEMENTS = 15
)

var (
	globalStyle = lipgloss.NewStyle().
			AlignVertical(lipgloss.Top).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#C8C093")).
			Width(40).
			Height(20).
			PaddingLeft(2)

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8C093")).
			PaddingLeft(2).
			PaddingBottom(1).
			PaddingTop(1)

	boldStyle = lipgloss.NewStyle().Bold(true)

	errorStyle = lipgloss.NewStyle().
			PaddingTop(1).
			MarginBottom(1).
			Foreground(lipgloss.Color("#C34043")).
			Bold(true)

	sessionItemStyle  = lipgloss.NewStyle().PaddingLeft(4)
	dirItemStyle      = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#727169"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF9E3B"))
)

type item struct {
	name    string
	iType   itemType
	current bool
}
type itemType int

const (
	DIR itemType = iota
	SESSION
)

func (i item) FilterValue() string { return i.name }

type items []item

func (i items) Len() int              { return len(i) }
func (i items) String(idx int) string { return i[idx].name }

type model struct {
	allItems      items
	filteredItems items

	textInput textinput.Model
	help      help.Model
	index     int
	offset    int

	base          string
	width         int
	height        int
	targetSession string
	err           error
}

func initialModel(base string, dirs []fs.DirEntry, sessions []string, current string) model {
	l := []item{}
	for _, s := range sessions {
		l = append(l, item{
			name:    s,
			iType:   SESSION,
			current: s == current,
		})
	}
	for _, d := range dirs {
		name := d.Name()
		if !slices.Contains(sessions, name) {
			l = append(l, item{
				name:  d.Name(),
				iType: DIR,
			})
		}
	}
	filteredItems := slices.Clone(l)

	ti := textinput.New()
	ti.Placeholder = "filter sessions / directories"
	ti.Cursor.SetMode(cursor.CursorStatic)
	ti.Focus()

	help := help.New()

	return model{
		base:          base,
		allItems:      l,
		filteredItems: filteredItems,
		textInput:     ti,
		help:          help,
	}
}

func (m *model) cursorUp() {
	if m.index > 0 {
		m.index--
	}
	if m.index-m.offset <= 0 && m.offset > 0 {
		m.offset--
	}
}

func (m *model) cursorDown() {
	if m.index < len(m.filteredItems)-1 {
		m.index++
	}
	if m.index-m.offset > MAX_ELEMENTS-2 && m.offset < len(m.filteredItems)-MAX_ELEMENTS {
		m.offset++
	}
}

func (m *model) SelectedItem() *item {
	if m.index >= 0 && m.index < len(m.filteredItems) {
		return &m.filteredItems[m.index]
	}
	return nil
}
func (m *model) resetFilter() {
	m.textInput.SetValue("")
	m.updateFilter()
}
func (m *model) updateFilter() {
	query := m.textInput.Value()
	m.filteredItems = []item{}
	if query == "" {
		for _, item := range m.allItems {
			m.filteredItems = append(m.filteredItems, item)
		}
		return
	}
	matches := fuzzy.FindFrom(query, m.allItems)
	for i, match := range matches {
		if i >= MAX_ELEMENTS {
			break
		}
		m.filteredItems = append(m.filteredItems, m.allItems[match.Index])
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
		globalStyle = lipgloss.NewStyle().
			MarginTop(m.height / 4).
			Inherit(globalStyle)
	case tea.KeyMsg:
		m.err = nil
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Down):
			m.cursorDown()
		case key.Matches(msg, keys.Up):
			m.cursorUp()
		case key.Matches(msg, keys.Reset):
			m.resetFilter()
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, keys.Confirm):
			it := m.SelectedItem()
			if it == nil {
				return m, nil
			}
			if it.current {
				return m, tea.Quit
			}
			if it.iType == DIR {
				if err := tmux.CreateSession(m.base, it.name); err != nil {
					m.err = err
					return m, nil
				}
			}
			m.targetSession = it.name
			return m, tea.Quit
		default:
			m.index = 0
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	m.updateFilter()
	return m, cmd
}

func (m model) View() string {
	base := ""
	base += filterStyle.Render(m.textInput.View()) + "\n"

	for i, item := range m.filteredItems {
		if i < m.offset {
			continue
		}
		if i >= m.offset+MAX_ELEMENTS {
			break
		}
		if i == m.index {
			base += selectedItemStyle.Render("* "+fmt.Sprintf(item.name)) + "\n"
		} else if item.iType == DIR {
			base += dirItemStyle.Render(fmt.Sprintf(item.name)) + "\n"
		} else if item.iType == SESSION {
			base += sessionItemStyle.Render(fmt.Sprintf(item.name)) + "\n"
		}
	}
	base += "\n"

	lstr := globalStyle.Render(base)

	if m.err != nil {
		lstr += "\n" + errorStyle.Render("Error: "+m.err.Error())
	} else {
		lstr += "\n" + errorStyle.Render("")
	}

	lstr += "\n" + m.help.View(keys)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, lstr)
}
