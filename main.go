package main

import (
	"time"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/application"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/category"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/color"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/gnome"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/ui"
)

var (
	uiResponse int // -1 is quit, -2 is back
	mainMenu = []string{"Manage Application Categories", "Create New Category Folder", "Delete Existing Category Folder", "Clean and Sort Category Files", "Apply Category Folders in Gnome", "Restore Default Layout in Gnome", "Quit"}
	okCancelOptions = []string{"Confirm", "Cancel"}
)

// Checks if at least one category folder exists and informs the user if they don't
func catFoldersExist() bool {
	if len(category.List) < 2 {
		ui.MessageWait("No category folders currently exist")
		return false
	}

	return true
}

// Application category management
func manageApplicationCategories() {
	var (
		catIndex int = 0
		newCatIndex int = 0
		appIndex int = 0
	)

	if !catFoldersExist() {
		return
	}

	for {
		uiResponse = ui.List("Select a " + color.Add("red", "category folder") + " to edit its " + color.Add("yellow", "applications"), category.GetNames(category.List), catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		appIndex = 0

		for {
			uiResponse = ui.List("Select an " + color.Add("red", "application") + " to change its " + color.Add("yellow", "category folder"), application.GetNames(category.List[catIndex].Applications), appIndex)

			if uiResponse == -1 {
				return
			} else if uiResponse == -2 {
				break
			}

			appIndex = uiResponse
			uiResponse = ui.List("Select a new " + color.Add("red", "category folder") + " for " + color.Add("yellow", category.List[catIndex].Applications[appIndex].Name), category.GetNames(category.List), newCatIndex);

			if uiResponse == -1 {
				return
			} else if uiResponse != -2 {
				newCatIndex = uiResponse
				category.ChangeAppCategory(category.List[catIndex].Applications[appIndex], catIndex, newCatIndex)
			}
		}
	}
}

// Create a new category
func createCategoryFolder() {
	newCategory := ui.Input("Enter a new category name")

	if newCategory != "" {
		status, err := category.Create(newCategory)

		if status {
			ui.MessageWait("Successfully created the " + color.Add("yellow", newCategory) + " category folder")
		} else {
			ui.MessageWait(err)
		}
	}
}

// Delete an existing category
func deleteCategoryFolder() {
	var catIndex int = 0

	if !catFoldersExist() {
		return
	}

	for {
		catNames := category.GetNames(category.GetListWithoutUncategorized())
		uiResponse = ui.List("Select a " + color.Add("red", "category folder") + " to " + color.Add("yellow", "delete"), catNames, catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		uiResponse = ui.List(color.Add("red", "Delete") + " the category folder " + color.Add("yellow", catNames[catIndex]) + "?", okCancelOptions, 0)

		if uiResponse == -1 {
			return
		} else if uiResponse != -2 && uiResponse != 1 {
			status, err := category.Delete(catNames[catIndex])

			if !status {
				ui.MessageWait(err)
			}
		}
	}
}

// Clean up category files by re-writing each of them
func cleanCategoryFiles() {
	if !catFoldersExist() {
		return
	}

	category.CleanFiles()
	ui.MessageWait("The applications in each category file have been cleaned and sorted")
}

// Apply configured categories to the gnome application list
func applyGnomeCategoryFolders() {
	if !catFoldersExist() {
		return
	}

	ui.Message("Clearing old category folders...")
	gnome.ClearCategories()
	time.Sleep(3 * time.Second)
	ui.Message("Applying categories...")
	gnome.ApplyCategories()
	ui.MessageWait("Successfully applied gnome category folders")
}

// Clear the gnome application list categories
func clearGnomeCategoryFolders() {
	gnome.ClearCategories()
	ui.MessageWait("Successfully cleared category folders")
}

// The main menu loop when running interactively
func mainMenuLoop() {
	var menuIndex = 0

	for {
		uiResponse = ui.List(color.Add("yellow", "Main Menu"), mainMenu, menuIndex);

		if uiResponse == -1 || uiResponse == -2 || uiResponse == len(mainMenu) - 1 {
			break
		}

		menuIndex = uiResponse

		switch uiResponse {
			case 0:
				manageApplicationCategories()

			case 1:
				createCategoryFolder()

			case 2:
				deleteCategoryFolder()

			case 3:
				cleanCategoryFiles()

			case 4:
				applyGnomeCategoryFolders()

			case 5:
				clearGnomeCategoryFolders()
		}
	}

	ui.Message("Quitting...")
}

// Main menu
func main() {
	mainMenuLoop()
}
