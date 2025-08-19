package ui

import (
	"fmt"
	"log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputTitle string
	inputValue string
)

type (
	inputErrMsg error
)

type inputModel struct {
	textInput textinput.Model
	err error
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
				case tea.KeyEnter:
					inputValue = m.textInput.Value()
					return m, tea.Quit

				case tea.KeyCtrlC, tea.KeyEsc:
					inputValue = ""
					return m, tea.Quit
			}

		// We handle errors just like any other message
		case inputErrMsg:
			m.err = msg
			return m, nil
		}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	lightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#616161"))
	darkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#494949"))

	return fmt.Sprintf(
		lipgloss.NewStyle().PaddingTop(1).PaddingLeft(3).PaddingBottom(1).Render(inputTitle + "\n\n%s\n\n%s"),
		m.textInput.View(),
		lightStyle.Render("enter") + darkStyle.Render(" select • ") + lightStyle.Render("esc") + darkStyle.Render(" quit"),
	) + "\n"
}

func inputInitialModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputModel{
		textInput: ti,
		err: nil,
	}
}

func Input(title string) (string) {
	inputTitle = title
	p := tea.NewProgram(inputInitialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return inputValue
}
