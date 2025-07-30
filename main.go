package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

var uiResponse int

func main() {
	uiResponse = ui.LoadList("Select an application", xdg.DisplayFiles);

	if (uiResponse != -1) {
		fmt.Println(xdg.DisplayFiles[uiResponse])
	}
}
