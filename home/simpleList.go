package home

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultWidth = 24
const listHeight = 12

var (
	containerStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(1, 0)
	inactiveContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("243")).Padding(1, 0)
	titleStyle             = lipgloss.NewStyle().Margin(0, 2).Padding(0, 2).Bold(true).Background(lipgloss.Color("62")).Foreground(lipgloss.Color("#ffffff"))
	itemStyle              = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle      = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle        = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle              = list.DefaultStyles().HelpStyle.Padding(1, 4, 0)
	quitTextStyle          = lipgloss.NewStyle().Margin(1, 4, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, l list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// str := fmt.Sprintf("%d. %s", index+1, i)
	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == l.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type SimpleList struct {
	focused  bool
	list     list.Model
	choice   string
	quitting bool
}

func (sl SimpleList) Focus() {
	sl.focused = true
}

func (sl SimpleList) Blur() {
	sl.focused = false
}

func (sl SimpleList) Focused() bool {
	return sl.focused
}

func (sl SimpleList) Init() tea.Cmd {
	return nil
}

func (sl SimpleList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sl.list.SetWidth(msg.Width)
		return sl, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := sl.list.SelectedItem().(item)
			if ok {
				sl.choice = string(i)
			}
			return sl, tea.Quit
		}
	}

	var cmd tea.Cmd
	sl.list, cmd = sl.list.Update(msg)
	return sl, cmd
}

func (sl SimpleList) View() string {
	if sl.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", sl.choice))
	}
	if sl.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}

	if sl.Focused() {
		return containerStyle.Render(sl.list.View())
	} else {
		return inactiveContainerStyle.Render(sl.list.View())
	}
}

func NewTestSimpleList() SimpleList {
	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
		item("Okonomiyaki"),
		item("Pasta"),
		item("Fillet Mignon"),
		item("Caviar"),
		item("Just Wine"),
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	sl := SimpleList{list: l, focused: true}

	return sl
}

func NewSimpleList() SimpleList {
	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	sl := SimpleList{list: l, focused: true}

	return sl
}
