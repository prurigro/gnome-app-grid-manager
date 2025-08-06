package env

import "os"

var (
	XdgDataHome string
	XdgDataDirs string
)

func init() {
	// Use the configured XDG_DATA_HOME or fall back on the default
	if os.Getenv("XDG_DATA_HOME") != "" {
		XdgDataHome = os.Getenv("XDG_DATA_HOME")
	} else {
		XdgDataHome = os.Getenv("HOME") + "/.local/share"
	}

	// Use the configured XDG_DATA_DIRS or fall back on the default
	if os.Getenv("XDG_DATA_DIRS") != "" {
		XdgDataDirs = os.Getenv("XDG_DATA_DIRS")
	} else {
		XdgDataDirs = "/usr/local/share:/usr/share"
	}
}
