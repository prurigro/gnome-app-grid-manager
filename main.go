package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/categories"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
)

var (
	uiResponse int // -1 is quit, -2 is back
	mainMenu = []string{"Categorize Applications", "Quit"}
)

func manageCategories() {
	var (
		catIndex int = 0
		appIndex int = 0
		newCatIndex int = 0
		appLauncher xdg.Launcher
	)

	for {
		uiResponse = ui.LoadList("Select a category to edit", categories.CatNames, catIndex);
		catIndex = uiResponse

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex := uiResponse
		category := categories.List[catIndex]
		appIndex = 0

		for {
			uiResponse = ui.LoadList("Select an application to categorize", xdg.GetAppNames(category.Applications), appIndex)

			if uiResponse == -1 {
				return
			} else if uiResponse == -2 {
				break
			} else {
				appIndex = uiResponse
				appLauncher = category.Applications[appIndex]
				uiResponse = ui.LoadList("Select a new category for " + appLauncher.AppName, categories.CatNames, catIndex);

				if uiResponse == -1 {
					return
				} else if uiResponse != -2 {
					newCatIndex = uiResponse
					fmt.Println(newCatIndex)
				}
			}
		}
	}
}

func main() {
	for uiResponse != -1 && uiResponse != len(mainMenu) - 1 {
		uiResponse = ui.LoadList("Main Menu", mainMenu, 0);

		switch uiResponse {
			case 0:
				manageCategories()
		}
	}

	fmt.Println("Quitting...")
}
