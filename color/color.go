package color

import (
	"git.darkcloud.ca/kevin/gnome-appcat-manager/cli"
)

var escapes = map[string]string{
	"gray": "\033[90m",
	"red": "\033[91m",
	"green": "\033[92m",
	"yellow": "\033[93m",
	"blue": "\033[94m",
	"violet": "\033[95m",
	"teal": "\033[96m",
	"white": "\033[97m",
	"bold": "\033[1m",
	"reset": "\033[0m",
}

func add(color string, text string) (string) {
	if cli.IsInteractive {
		return escapes["bold"] + escapes[color] + text + escapes["reset"]
	} else {
		return text
	}
}

func Gray(text string) (string) {
	return add("gray", text)
}

func Red(text string) (string) {
	return add("red", text)
}

func Green(text string) (string) {
	return add("green", text)
}

func Yellow(text string) (string) {
	return add("yellow", text)
}

func Blue(text string) (string) {
	return add("blue", text)
}

func Violet(text string) (string) {
	return add("violet", text)
}

func Teal(text string) (string) {
	return add("teal", text)
}

func White(text string) (string) {
	return add("white", text)
}
