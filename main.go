package main

import (
	"fmt"
	"os"
	"time"
	"github.com/prurigro/gnome-app-grid-manager/application"
	"github.com/prurigro/gnome-app-grid-manager/category"
	"github.com/prurigro/gnome-app-grid-manager/color"
	"github.com/prurigro/gnome-app-grid-manager/gnome"
	"github.com/prurigro/gnome-app-grid-manager/ui"
)

var (
	appName string
	appVersion string = "v1.0.1"
	mainMenuOptions = []string{"Manage application categories", "Create new category folder", "Delete existing category folder", "Apply category folders in Gnome", "Restore default layout in Gnome", "Clean and sort data files", "Quit"}
	okCancelOptions = []string{"Confirm", "Cancel"}
	uiResponse int // -1 is quit, -2 is back
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
		uiResponse = ui.List("Select a " + color.Red("category folder") + " to edit its " + color.Yellow("applications"), category.GetNames(category.List), catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		appIndex = 0

		for {
			uiResponse = ui.List("Select an " + color.Red("application") + " to change its " + color.Yellow("category folder"), application.GetNames(category.List[catIndex].Applications), appIndex)

			if uiResponse == -1 {
				return
			} else if uiResponse == -2 {
				break
			}

			appIndex = uiResponse
			uiResponse = ui.List("Select a new " + color.Red("category folder") + " for " + color.Yellow(category.List[catIndex].Applications[appIndex].Name), category.GetNames(category.List), newCatIndex);

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
	newCategory := ui.Input("Enter a new " + color.Red("category folder") + " name")

	if newCategory != "" {
		status, err := category.Create(newCategory)

		if status {
			ui.MessageWait("Successfully created the " + color.Red(newCategory) + " category folder")
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
		uiResponse = ui.List("Select a " + color.Red("category folder") + " to " + color.Yellow("delete"), catNames, catIndex);

		if uiResponse == -1 || uiResponse == -2 {
			return
		}

		catIndex = uiResponse
		uiResponse = ui.List(color.Red("Delete") + " the category folder " + color.Yellow(catNames[catIndex]) + "?", okCancelOptions, 0)

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
	ui.MessageWait("The " + color.Red("applications") + " in each " + color.Yellow("category file") + " have been cleaned and sorted")
}

// Clear the gnome application list categories
func restoreGnomeDefaultLayout() {
	gnome.RestoreDefault()
	ui.MessageWait("Successfully restored the " + color.Red("default layout") + " in " + color.Yellow("Gnome"))
}

// Apply configured categories to the gnome application list
func applyGnomeCategoryFolders() {
	if !catFoldersExist() {
		return
	}

	ui.Message("Removing the existing " + color.Red("category folders") + " from " + color.Yellow("Gnome") + "...")
	gnome.RestoreDefault()
	time.Sleep(3 * time.Second)
	ui.Message("Applying configured " + color.Red("category folders") + "...")
	gnome.ApplyCategories()
	ui.MessageWait("Successfully applied " + color.Red("category folders") + " in " + color.Yellow("Gnome"))
}

// The main menu loop when running interactively
func mainMenuLoop() {
	var menuIndex = 0

	for {
		uiResponse = ui.List(color.Yellow("Main Menu"), mainMenuOptions, menuIndex);

		if uiResponse == -1 || uiResponse == -2 || uiResponse == len(mainMenuOptions) - 1 {
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
				applyGnomeCategoryFolders()

			case 4:
				restoreGnomeDefaultLayout()

			case 5:
				cleanCategoryFiles()
		}
	}

	ui.Message("Quitting...")
}

// Show version information and exit
func displayVersion() {
	fmt.Println("\n" + color.White(appName) + " " + color.Violet(appVersion) + "\n")
	os.Exit(0)
}

// Show the help text and exit
func displayHelp(status int) {
	fmt.Println("\n" + color.Yellow(appName) + " - Organize your Gnome overview applications by category")
	fmt.Println("\n" + color.Blue("USAGE"))
	fmt.Println("  " + color.Violet(appName) + "\t\tRun interactively")
	fmt.Println("  " + color.Violet(appName) + " " + color.Gray("[") + color.White("option") + color.Gray("]") + "\tDirectly run one of the options below")
	fmt.Println("\n" + color.Blue("OPTIONS"))
	fmt.Println("  " + color.White("-a") + color.Gray("|") + color.White("--apply") + "\t" + mainMenuOptions[3])
	fmt.Println("  " + color.White("-r") + color.Gray("|") + color.White("--restore") + "\t" + mainMenuOptions[4])
	fmt.Println("  " + color.White("-c") + color.Gray("|") + color.White("--clean") + "\t" + mainMenuOptions[5])
	fmt.Println("  " + color.White("-v") + color.Gray("|") + color.White("--version") + "\t" + "Display the current version")
	fmt.Println("  " + color.White("-h") + color.Gray("|") + color.White("--help") + "\t" + "Show this help text")
	fmt.Println("\n" + color.Blue("DATA DIRECTORY"))
	fmt.Println("  " + category.Directory)
	fmt.Println("")
	os.Exit(status)
}

// Main menu
func main() {
	// Store the app name
	appName = os.Args[0]

	// Get the command line arguments
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("\n" + color.Red("Error") + ": Multiple arguments are not supported")
		displayHelp(1)
	} else if len(args) == 1 {
		switch args[0] {
			case "-a", "--apply":
				applyGnomeCategoryFolders()

			case "-r", "--restore":
				restoreGnomeDefaultLayout()

			case "-c", "--clean":
				cleanCategoryFiles()

			case "-v", "--version":
				displayVersion()

			case "-h", "--help":
				displayHelp(0)

			default:
				fmt.Println(color.Red("Error") + ": Invalid argument")
				displayHelp(1)
		}
	} else {
		// Run the main menu loop if no arguments are provided
		mainMenuLoop()
	}
}
