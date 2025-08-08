package cli

import (
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

var IsInteractive bool

func init() {
	// Check if the terminal session is interactive
	IsInteractive = terminal.IsTerminal(int(os.Stdout.Fd()))
}
