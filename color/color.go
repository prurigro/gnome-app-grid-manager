package color

import (
	"git.darkcloud.ca/kevin/gnome-appcat-manager/cli"
)

var escapes = map[string]string{
	"grey": "\033[90m",
	"red": "\033[91m",
	"green": "\033[92m",
	"yellow": "\033[93m",
	"blue": "\033[94m",
	"violet": "\033[95m",
	"teal": "\033[96m",
	"white": "\033[97m",
	"reset": "\033[0m",
}

func Add(color string, text string) (string) {
	if cli.IsInteractive() {
		return escapes[color] + text + escapes["reset"]
	} else {
		return text
	}
}
