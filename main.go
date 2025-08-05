package main

import (
	"git.darkcloud.ca/kevin/gnome-appcat-manager/application"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/category"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/color"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

var (
	uiResponse int // -1 is quit, -2 is back
	mainMenu = []string{"Manage Categories", "Create Category", "Delete Category", "Quit"}
	okCancelOptions = []string{"Confirm", "Cancel"}
)

// Interactive application categorization
func manageCategories() {
	var (
		catIndex int = 0
		newCatIndex int = 0
		appIndex int = 0
	)

	for {
		uiResponse = ui.List("Select a " + color.Add("red", "category") + " to edit its " + color.Add("yellow", "applications"), category.GetNames(category.List), catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		appIndex = 0

		for {
			uiResponse = ui.List("Select an " + color.Add("red", "application") + " to change its " + color.Add("yellow", "category"), application.GetNames(category.List[catIndex].Applications), appIndex)

			if uiResponse == -1 {
				return
			} else if uiResponse == -2 {
				break
			}

			appIndex = uiResponse
			uiResponse = ui.List("Select a new " + color.Add("red", "category") + " for " + color.Add("yellow", category.List[catIndex].Applications[appIndex].Name), category.GetNames(category.List), catIndex);

			if uiResponse == -1 {
				return
			} else if uiResponse != -2 {
				newCatIndex = uiResponse
				category.ChangeAppCategory(category.List[catIndex].Applications[appIndex], catIndex, newCatIndex)
			}
		}
	}
}

func createCategory() {

}

func deleteCategory() {
	var catIndex int = 0

	for {
		catNames := category.GetNames(category.GetListWithoutUncategorized())
		uiResponse = ui.List("Select a " + color.Add("red", "category") + " to " + color.Add("yellow", "delete"), catNames, catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		uiResponse = ui.List(color.Add("red", "Delete") + " the category " + color.Add("yellow", catNames[catIndex]) + "?", okCancelOptions, 0)

		if uiResponse == -1 {
			return
		} else if uiResponse != -2 && uiResponse != 1 {
			status, err := category.Delete(catNames[catIndex])

			if !status {
				ui.Message(err, true)
			}
		}
	}
}

func main() {
	for uiResponse != -1 && uiResponse != len(mainMenu) - 1 {
		uiResponse = ui.List(color.Add("yellow", "Main Menu"), mainMenu, 0);

		switch uiResponse {
			case 0:
				manageCategories()

			case 1:
				createCategory()

			case 2:
				deleteCategory()
		}
	}

	ui.Message("Quitting...", false)
}
