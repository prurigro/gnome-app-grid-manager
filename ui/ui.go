package ui

import (
	"bufio"
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
var selectedValue int = -1

// Style
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(1)
	noItemsStyle      = lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("8"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(3)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("12"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(3)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(3).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 2)
)

// Additional keybindings for help
type keyMap struct {
	Enter key.Binding
	Backspace key.Binding
}

func (k keyMap) AdditionalShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Backspace}
}

func (k keyMap) AdditionalFullHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Backspace}
}

var keys = keyMap{
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
type item struct {
	index int
	value string
}

func (i item) FilterValue() string {
	return ""
}

// Array of items
var items = []list.Item{}

// List item format
type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)

	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.value)
	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					selectedValue = -2
					return m, tea.Quit

				case "enter":
					i, ok := m.list.SelectedItem().(item)

					if ok {
						selectedValue = i.index
						m.choice = i.value
					}

					return m, tea.Quit
			}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return ""
	}

	if m.quitting {
		return ""
	}

	return "\n" + m.list.View()
}

func List(title string, listItems []string, startingIndex int) (int) {
	selectedValue = -1
	items = nil

	for index, value := range listItems {
		items = append(items, item{index, value})
	}

	l := list.New(items, itemDelegate{}, 1, 1)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.NoItems = noItemsStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.AdditionalShortHelpKeys = keys.AdditionalShortHelp
	l.AdditionalFullHelpKeys = keys.AdditionalFullHelp

	if startingIndex >= len(listItems) - 1 {
		startingIndex = len(listItems) - 1
	}

	l.Select(startingIndex)

	m := model{list: l}
	_, err := tea.NewProgram(m).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return selectedValue
}

func Message(message string, waitForEnter bool) {
	// Display the message
	fmt.Println(lipgloss.NewStyle().PaddingTop(1).PaddingLeft(3).PaddingBottom(1).Render(message))

	if waitForEnter {
		// Inform the user of what to do next
		fmt.Println(lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("8")).Render("Press enter to continue"))

		// Hide the cursor
		fmt.Print("\033[?25l")

		// Show the cursor when this function concludes
		defer fmt.Print("\033[?25h")

		// Create a reader and wait for the enter key
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}
}
