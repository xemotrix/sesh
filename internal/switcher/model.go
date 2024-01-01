package switcher

import (
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"

	"github.com/xemotrix/sesh/internal/tmux"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8C093")).
			PaddingLeft(2).
			PaddingBottom(1).
			PaddingTop(1)

	boldStyle = lipgloss.NewStyle().Bold(true)

	errorStyle = lipgloss.NewStyle().
			PaddingTop(1).
			Foreground(lipgloss.Color("#C34043")).
			Bold(true)

	sessionItemStyle  = lipgloss.NewStyle().PaddingLeft(4)
	dirItemStyle      = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#727169"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF9E3B"))
)

type model struct {
	list          list.Model
	base          string
	width         int
	height        int
	targetSession string
	err           error
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := i.name

	if i.current {
		str = boldStyle.Render(str)
	}

	var fn func(strs ...string) string
	if i.iType == SESSION {
		fn = sessionItemStyle.Render
	} else {
		fn = dirItemStyle.Render
	}

	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("* " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

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

func initialModel(base string, dirs []fs.DirEntry, sessions []string, current string) model {
	l := []list.Item{}
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
	li := list.New(l, itemDelegate{}, 20, 20)
	li.SetShowTitle(false)
	li.SetShowStatusBar(false)
	li.SetShowFilter(false)
	ti := textinput.New()
	ti.Placeholder = "filter sessions / directories"
	ti.Cursor.SetMode(cursor.CursorStatic)
	li.FilterInput = ti

	li.SetShowHelp(false)
	li.SetShowPagination(false)

	return model{
		base: base,
		list: li,
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
		case "ctrl+j":
			m.list.CursorDown()
		case "ctrl+k":
			m.list.CursorUp()
		case "enter":
			if m.list.SelectedItem() == nil {
				return m, nil
			}
			it := (m.list.SelectedItem()).(item)
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
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	base := ""
	if m.list.SettingFilter() {
		base += filterStyle.Render(m.list.FilterInput.View()) + "\n"
	} else {
		base += "\n"
	}
	base += m.list.View()
	base += "\n"

	lstr := globalStyle.Render(base)

	if m.err != nil {
		lstr += "\n" + errorStyle.Render("Error: "+m.err.Error())
	} else {
		lstr += "\n"
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lstr)
}
