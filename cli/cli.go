package cli

import (
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

// Check if the terminal session is interactive
func IsInteractive() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
