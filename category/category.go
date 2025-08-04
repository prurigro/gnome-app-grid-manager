package category

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/application"
)

type Data struct {
	File string
	Name string
	Applications []application.Data
}

var (
	List = []Data{}
	Names = []string{}
	Files = []string{}
	categoriesDirectory = os.Getenv("XDG_DATA_HOME") + "/gnome-shell/categories"
)

// Creates the categoriesDirectory if it doesn't already exist
func createCatDirWhenMissing() {
	if stat, err := os.Stat(categoriesDirectory); err != nil || !stat.IsDir() {
		os.MkdirAll(categoriesDirectory, 0755)
	}
}

// Move an application from one category to another
func ChangeAppCategory(appItem application.Data, oldCatIndex int, newCatIndex int) {
	fmt.Println("Moving " + appItem.Name + " from " + List[oldCatIndex].Name + " to " + List[newCatIndex].Name)
	os.Exit(0)
}

// Updates the list of all file names
func UpdateFiles() {
	var newFiles []string

	for _, item := range List {
		newFiles = append(newFiles, item.File)
	}

	Files = newFiles
}

// Updates the list of all category names
func UpdateNames() {
	var newNames []string

	for _, item := range List {
		newNames = append(newNames, item.Name)
	}

	Names = newNames
}

// Populate the List of categories
func Populate() {
	var files []application.Data

	// Create the categories directory if it's missing
	createCatDirWhenMissing()

	// Regex to pull just the name out of the category file
	catNameRegex := regexp.MustCompile("(?i)(.*)\\.category")

	// Initialize the List with all the applications in uncategorized
	List = []Data{ { File: "", Name: "Uncategorized", Applications: []application.Data{} } }

	// Store the application file names and List so we can track and populate uncategorized items
	appFiles := application.Files
	appList := application.List

	// Read the files in the categories directory
	catFiles, _ := os.ReadDir(categoriesDirectory)

	for _, catFile := range catFiles {
		if catNameRegex.MatchString(catFile.Name()) {
			files = nil
			catNameMatchReference := catNameRegex.FindStringSubmatch(catFile.Name())
			filePath := categoriesDirectory + "/" + catFile.Name()
			file, err := os.Open(filePath)

			if err != nil {
				log.Fatal("Unable to open the file " + filePath)
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

			// Sort alphabetically (case insensitive)
			sort.Slice(files, func(x, y int) bool {
				return strings.ToLower(files[x].Name) < strings.ToLower(files[y].Name)
			})

			List = append(List, Data{
				File: catFile.Name(),
				Name: catNameMatchReference[1],
				Applications: files,
			})
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

	UpdateFiles()
	UpdateNames()
}

func init() {
	Populate()
}
