package main

import (
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

func main() {
	ui.LoadList("Select an application", xdg.DisplayFiles);
}
