package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prurigro/gnome-app-grid-manager/cli"
)

type messageModel struct {
	quit bool
}

func (m messageModel) Init() tea.Cmd {
	return nil
}

func (m messageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
				case tea.KeyEnter, tea.KeyEsc, tea.KeyCtrlC:
					return m, tea.Quit
			}
	}

	return m, nil
}

func (m messageModel) View() string {
	if m.quit {
		return ""
	}

	// Inform the user of how to continue
	return lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("8")).Render("Press enter or esc to continue")
}

// Clear the terminal and display a message
func Message(message string) {
	if cli.IsInteractive {
		fmt.Print("\033[H\033[2J")
	}

	fmt.Println(lipgloss.NewStyle().PaddingTop(1).PaddingLeft(3).PaddingBottom(1).Render(message))
}

// Display a message and wait for them to press enter
func MessageWait(message string) {
	Message(message)

	if cli.IsInteractive {
		// Listen for enter or escape
		p := tea.NewProgram(messageModel{})
		p.Start()
	}
}
