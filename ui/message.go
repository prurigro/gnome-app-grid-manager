package ui

import (
	"bufio"
	"fmt"
	"os"
	"github.com/charmbracelet/lipgloss"
	"git.darkcloud.ca/kevin/gnome-app-grid-manager/cli"
)

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
