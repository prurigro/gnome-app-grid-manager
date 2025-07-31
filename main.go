package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/categories"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

var uiResponse int

func main() {
	uiResponse = ui.LoadList("Main Menu", []string{"List Applications", "List Categories", "Quit"});

	if (uiResponse == 0) {
		uiResponse = ui.LoadList("Select an application", xdg.AppNames());

		if (uiResponse != -1) {
			fmt.Println(xdg.FileNames()[uiResponse])
		}
	} else if (uiResponse == 1) {
		uiResponse = ui.LoadList("Select a category", categories.CatNames());

		if (uiResponse != -1) {
			fmt.Println(categories.FileNames()[uiResponse])
		}
	} else {
		fmt.Println("Quitting...")
	}
}
