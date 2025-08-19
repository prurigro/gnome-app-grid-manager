package category

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"github.com/prurigro/gnome-app-grid-manager/application"
	"github.com/prurigro/gnome-app-grid-manager/color"
	"github.com/prurigro/gnome-app-grid-manager/env"
)

type Data struct {
	File string
	Name string
	Applications []application.Data
}

var (
	List = []Data{}
	Directory = env.XdgDataHome + "/gnome-shell/categories"
)

// Write a category to file
func writeCategory(catItem Data) {
	// Only write if a file exists (ie: don't write uncategorized)
	if catItem.File != "" {
		// The full category file path
		filePath := Directory + "/" + catItem.File

		// Truncate and create the category file
		file, err := os.Create(filePath)

		// Fail with an error if we can't open the file
		if  err != nil {
			log.Fatal("Unable to open the file " + color.Red(filePath))
		}

		// Close the file once we're finished here
		defer file.Close()

		// Write the list of applications to the category file
		writer := bufio.NewWriter(file)

		// Write the applications
		for index, appItem := range catItem.Applications {
			// write the application filename
			_, err := writer.WriteString(appItem.File)

			if err != nil {
				panic(err)
			}

			// Add a line break if this wasn't the last file
			if index < len(catItem.Applications) - 1 {
				_, err := writer.WriteString("\n")

				if err != nil {
					panic(err)
				}
			}
		}

		if err := writer.Flush(); err != nil {
			log.Fatal("Unable to write to the file " + color.Red(filePath))
		}
	}
}

// Remove an application from a category
func removeApplication(appItem application.Data, catIndex int) {
	var newApplications = []application.Data{}

	for _, item := range List[catIndex].Applications {
		if item.File != appItem.File || item.Name != appItem.Name {
			newApplications = append(newApplications, item)
		}
	}

	List[catIndex].Applications = newApplications
	writeCategory(List[catIndex])
}

// Add an application to a category
func addApplication(appItem application.Data, catIndex int) {
	List[catIndex].Applications = append(List[catIndex].Applications, appItem)
	application.Sort(&List[catIndex].Applications)
	writeCategory(List[catIndex])
}

// Clean up category files by re-writing each of them
func CleanFiles() {
	for _, catItem := range GetListWithoutUncategorized() {
		writeCategory(catItem)
	}
}

// Move an application from one category to another
func ChangeAppCategory(appItem application.Data, oldCatIndex int, newCatIndex int) {
	removeApplication(appItem, oldCatIndex)
	addApplication(appItem, newCatIndex)
}

// Create a category
func Create(name string) (bool, string) {
	// Complain if the category already exists
	if slices.Contains(GetNames(GetListWithoutUncategorized()), name) {
		return false, "The category " + color.Red(name) + " already exists"
	}

	// The file name and the full file path
	fileName := name + ".category"
	filePath := Directory + "/" + fileName

	// Create the categories directory if it doesn't already exist
	if stat, err := os.Stat(Directory); err != nil || !stat.IsDir() {
		os.MkdirAll(Directory, 0755)
	}

	// Create the category file if it doesn't exist (otherwise our job here is already done)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the category file and complain if it fails
		file, err := os.Create(filePath)

		if err != nil {
			log.Fatal("The file " + color.Red(fileName) + " could not be created")
		}

		// Close the file
		file.Close()
	}

	// Re-populate List
	Populate()

	// Return successfully
	return true, ""
}

// Delete a category
func Delete(name string) (bool, string) {
	// The file name and the full file path
	fileName := name + ".category"
	filePath := Directory + "/" + fileName

	// Delete the category file if it exists (otherwise our job here is already done)
	if _, err := os.Stat(filePath); err == nil || !os.IsNotExist(err) {
		// Delete the category file and complain if it fails
		err := os.Remove(filePath)

		if err != nil {
			log.Fatal("The file " + color.Red(fileName) + " could not be deleted")
		}
	}

	// Re-populate List
	Populate()

	// Return successfully
	return true, ""
}

// Rename a category
func Rename(catIndex int, newName string) (bool, string) {
	// Complain if the new category already exists
	if slices.Contains(GetNames(GetListWithoutUncategorized()), newName) {
		return false, "The category " + color.Red(newName) + " already exists"
	}

	// Store the old category name
	oldName := List[catIndex].Name

	// Rename the category in List
	List[catIndex].Name = newName
	List[catIndex].File = newName + ".category"

	// Create a new category file from the updated category data
	writeCategory(List[catIndex])

	// Delete the old category file and re-populate
	return Delete(oldName)
}

// Retrieve the list of categories without uncategorized
func GetListWithoutUncategorized() ([]Data) {
	var listWithoutUncategorized []Data

	for _, item := range List {
		if item.File != "" {
			listWithoutUncategorized = append(listWithoutUncategorized, item)
		}
	}

	return listWithoutUncategorized
}

// Returns an array of file names from a list of categories
func GetFiles(categories []Data) []string {
	var files []string

	for _, item := range categories {
		files = append(files, item.File)
	}

	return files
}

// Returns an array of app names from a list of categories
func GetNames(categories []Data) []string {
	var names []string

	for _, item := range categories {
		names = append(names, item.Name)
	}

	return names
}

// Populate the List of categories
func Populate() {
	var files []application.Data

	// Regex to pull just the name out of the category file
	catNameRegex := regexp.MustCompile("(?i)(.*)\\.category")

	// (Re-)Initialize the List with all the applications in uncategorized
	List = []Data{ { File: "", Name: "Uncategorized", Applications: []application.Data{} } }

	// Store the application file names and List so we can track and populate uncategorized items
	appFiles := application.GetFiles(application.List)
	appList := make([]application.Data, len(application.List))
	copy(appList, application.List[:])

	// The set of files in Directory
	catFiles, _ := os.ReadDir(Directory)

	// Loop through each file in Directory if it exists
	if stat, err := os.Stat(Directory); err == nil && stat.IsDir() {
		for _, catFile := range catFiles {
			// Only open .category files
			if catNameRegex.MatchString(catFile.Name()) {
				// Empty out the files array so we can start adding a fresh set from the latest category
				files = nil

				// Match the category file name so we can pull the category name out from it later
				catNameMatchReference := catNameRegex.FindStringSubmatch(catFile.Name())

				// Get the full path of the file
				filePath := Directory + "/" + catFile.Name()

				// Open the file
				file, err := os.Open(filePath)

				// Fail with an error if we can't open the file
				if err != nil {
					log.Fatal("Unable to open the file " + color.Red(filePath))
				}

				// Read the category file and build the list of applications for that category
				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					if application.FileMatchRegex.MatchString(scanner.Text()) {
						appFile := application.FileMatchRegex.FindStringSubmatch(scanner.Text())[0]

						if (slices.Contains(appFiles, appFile)) {
							index := slices.Index(appFiles, appFile)
							files = append(files, appList[index])
							_ = slices.Delete(appFiles, index, index + 1)
							_ = slices.Delete(appList, index, index + 1)
						}
					}
				}

				// Close the file
				file.Close()

				// Sort alphabetically (case insensitive)
				sort.Slice(files, func(x, y int) bool {
					return strings.ToLower(files[x].Name) < strings.ToLower(files[y].Name)
				})

				// Add the category to List
				List = append(List, Data{
					File: catFile.Name(),
					Name: catNameMatchReference[1],
					Applications: files,
				})
			}
		}
	}

	// Populate the uncategorized list using the items remaining in the appList
	for _, item := range appList {
		if item.File != "" {
			List[0].Applications = append(List[0].Applications, item)
		}
	}

	// Sort alphabetically (case insensitive)
	sort.Slice(List, func(x, y int) bool {
		return strings.ToLower(List[x].File) < strings.ToLower(List[y].File)
	})
}

func init() {
	Populate()
}
