package switcher

import (
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"

	"github.com/xemotrix/sesh/tmux"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	globalStyle = lipgloss.NewStyle().
			AlignVertical(lipgloss.Bottom).
			Width(50).
			PaddingLeft(2)

	borderedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("65"))

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8C093")).
			PaddingLeft(2)

	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
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
	name  string
	iType itemType
}
type itemType int

const (
	DIR itemType = iota
	SESSION
)

func (i item) FilterValue() string { return i.name }

func initialModel(base string, dirs []fs.DirEntry, sessions []string) model {
	l := []list.Item{}
	for _, s := range sessions {
		l = append(l, item{
			name:  s,
			iType: SESSION,
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
			it := (m.list.SelectedItem()).(item)
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
	if m.err != nil {
		base += "\n" + errorStyle.Render("Error: "+m.err.Error())
	}

	lstr := globalStyle.Render(base)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lstr)
}
