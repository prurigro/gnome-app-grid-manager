package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/category"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/application"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

var (
	uiResponse int // -1 is quit, -2 is back
	mainMenu = []string{"Categorize Applications", "Quit"}
)

// Interactive application categorization
func categorizeApplications() {
	var (
		catIndex int = 0
		newCatIndex int = 0
		appIndex int = 0
		appItem application.Data
	)

	for {
		uiResponse = ui.LoadList("Select a category to edit", category.Names, catIndex);
		catIndex = uiResponse

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex := uiResponse
		catItem := category.List[catIndex]
		appIndex = 0

		for {
			uiResponse = ui.LoadList("Select an application to categorize", application.GetNames(catItem.Applications), appIndex)

			if uiResponse == -1 {
				return
			} else if uiResponse == -2 {
				break
			} else {
				appIndex = uiResponse
				appItem = catItem.Applications[appIndex]
				uiResponse = ui.LoadList("Select a new category for " + appItem.Name, category.Names, catIndex);

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
				categorizeApplications()
		}
	}

	fmt.Println("Quitting...")
}
