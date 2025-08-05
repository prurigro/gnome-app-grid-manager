package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// Selected value
var listSelectedValue int = -1

// Style
var (
	listTitleStyle        = lipgloss.NewStyle().MarginLeft(1)
	listNoItemsStyle      = lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("8"))
	listItemStyle         = lipgloss.NewStyle().PaddingLeft(3)
	listSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("12"))
	listPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(3)
	listHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(3).PaddingBottom(1)
	listQuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 2)
)

// Additional keybindings for help
type listKeyMap struct {
	Enter key.Binding
	Backspace key.Binding
}

func (k listKeyMap) AdditionalShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Backspace}
}

func (k listKeyMap) AdditionalFullHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Backspace}
}

var listKeys = listKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("bksp", "back"),
	),
}

// List item definition
type listItem struct {
	index int
	value string
}

func (i listItem) FilterValue() string {
	return ""
}

// Array of items
var listItems = []list.Item{}

// List item format
type listItemDelegate struct{}

func (d listItemDelegate) Height() int {
	return 1
}

func (d listItemDelegate) Spacing() int {
	return 0
}

func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)

	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.value)
	fn := listItemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return listSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type listModel struct {
	list list.Model
	choice string
	quitting bool
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetSize(msg.Width, msg.Height - 1)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
				case "q", "ctrl+c":
					m.quitting = true
					return m, tea.Quit

				case "backspace":
					listSelectedValue = -2
					return m, tea.Quit

				case "enter":
					i, ok := m.list.SelectedItem().(listItem)

					if ok {
						listSelectedValue = i.index
						m.choice = i.value
					}

					return m, tea.Quit
			}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.choice != "" {
		return ""
	}

	if m.quitting {
		return ""
	}

	return "\n" + m.list.View()
}

func List(title string, items []string, startingIndex int) (int) {
	listSelectedValue = -1
	listItems = nil

	for index, value := range items {
		listItems = append(listItems, listItem{index, value})
	}

	l := list.New(listItems, listItemDelegate{}, 1, 1)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.NoItems = listNoItemsStyle
	l.Styles.PaginationStyle = listPaginationStyle
	l.Styles.HelpStyle = listHelpStyle

	l.AdditionalShortHelpKeys = listKeys.AdditionalShortHelp
	l.AdditionalFullHelpKeys = listKeys.AdditionalFullHelp

	if startingIndex >= len(items) - 1 {
		startingIndex = len(items) - 1
	}

	l.Select(startingIndex)

	m := listModel{list: l}
	_, err := tea.NewProgram(m).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return listSelectedValue
}
