package gnome

import (
	"encoding/json"
	"log"
	"os/exec"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/application"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/category"
)

func RestoreDefault() () {
	var (
		err error
		cmd *exec.Cmd
	)

	// Reset the category list
	cmd = exec.Command("gsettings", "set", "org.gnome.desktop.app-folders", "folder-children", "[]")
	_, err = cmd.Output()

	if err != nil {
		log.Fatal("Error resetting the category list in dconf:", err)
	}

	// Reset the category folders
	cmd = exec.Command("dconf", "reset", "-f", "/org/gnome/desktop/app-folders/folders/")
	_, err = cmd.Output()

	if err != nil {
		log.Fatal("Error resetting the category folders in dconf:", err)
	}
}

func ApplyCategories() {
	var (
		categories []category.Data = category.GetListWithoutUncategorized()
		cmd *exec.Cmd
		err error
		jsonList []byte
	)

	for _, catItem := range categories {
		// Create the category folder
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.app-folders.folder:/org/gnome/desktop/app-folders/folders/" + catItem.Name + "/", "name", catItem.Name)
		_, err = cmd.Output()

		if err != nil {
			log.Fatal("Error creating the category folder in dconf:", err)
		}

		// Add applications to the category if any exist
		if len(catItem.Applications) > 0 {
			jsonList, err = json.Marshal(application.GetFiles(catItem.Applications))

			if err != nil {
				log.Fatal("Error converting " + catItem.Name + " application names to json:", err)
			}

			cmd = exec.Command("gsettings", "set", "org.gnome.desktop.app-folders.folder:/org/gnome/desktop/app-folders/folders/" + catItem.Name + "/", "apps", string(jsonList))
			_, err = cmd.Output()

			if err != nil {
				log.Fatal("Error adding " + catItem.Name + " applications in dconf:", err)
			}
		}
	}

	// Build the category list
	jsonList, err = json.Marshal(category.GetNames(categories))

	if err != nil {
		log.Fatal("Error converting category names to json:", err)
	}

	// Apply the category list
	cmd = exec.Command("gsettings", "set", "org.gnome.desktop.app-folders", "folder-children", string(jsonList))
	_, err = cmd.Output()

	if err != nil {
		log.Fatal("Error applying the category list in dconf:", err)
	}

	// Reset the app picker layout
	cmd = exec.Command("gsettings", "set", "org.gnome.shell", "app-picker-layout", "[]")
	_, err = cmd.Output()

	if err != nil {
		log.Fatal("Error resetting the app picker layout in dconf:", err)
	}
}
