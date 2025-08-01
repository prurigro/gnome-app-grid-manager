package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/categories"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
)

var (
	uiResponse int // -1 is quit, -2 is back
	mainMenu = []string{"Categories", "Quit"}
)

func manageCategories() {
	for {
		uiResponse = ui.LoadList("Select a category", categories.CatNames);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex := uiResponse
		category := categories.List[catIndex]
		uiResponse = ui.LoadList("Select an application", xdg.GetAppNames(category.Applications))

		if uiResponse == -1 {
			return
		}
	}
}

func main() {
	for uiResponse != -1 && uiResponse != len(mainMenu) - 1 {
		uiResponse = ui.LoadList("Main Menu", mainMenu);

		switch uiResponse {
			case 0:
				manageCategories()
		}
	}

	fmt.Println("Quitting...")
}
