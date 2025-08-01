package categories

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
)

type category struct {
	File string
	Name string
	Applications []xdg.Launcher
}

var (
	List = []category{}
	CatNames = []string{}
	FileNames = []string{}
	categoriesDirectory = os.Getenv("XDG_DATA_HOME") + "/gnome-shell/categories"
)

func createCatDirWhenMissing() {
	if stat, err := os.Stat(categoriesDirectory); err != nil || !stat.IsDir() {
		os.MkdirAll(categoriesDirectory, 0755)
	}
}

// Updates the list of all file names
func UpdateFileNames() {
	var newFileNames []string

	for _, item := range List {
		newFileNames = append(newFileNames, item.File)
	}

	FileNames = newFileNames
}

// Updates the list of all category names
func UpdateCatNames() {
	var newCatNames []string

	for _, item := range List {
		newCatNames = append(newCatNames, item.Name)
	}

	CatNames = newCatNames
}

// Populate the List of categories
func Populate() {
	var files []xdg.Launcher

	// Create the categories directory if it's missing
	createCatDirWhenMissing()

	// Regex to pull just the name out of the category file
	catNameRegex := regexp.MustCompile("(?i)(.*)\\.category")

	// Initialize the List with all the applications in uncategorized
	List = []category{ { File: "", Name: "Uncategorized", Applications: []xdg.Launcher{} } }

	// Store the xdg filesnames and List so we can track and populate uncategorized items
	xdgFileNames := xdg.FileNames
	xdgList := xdg.List

	// Read the files in the categories directory
	categoryFiles, _ := os.ReadDir(categoriesDirectory)

	for _, catFile := range categoryFiles {
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
				if xdg.FileMatchRegex.MatchString(scanner.Text()) {
					filename := xdg.FileMatchRegex.FindStringSubmatch(scanner.Text())[0]

					if (slices.Contains(xdgFileNames, filename)) {
						index := slices.Index(xdgFileNames, filename)
						files = append(files, xdgList[index])
						_ = slices.Delete(xdgFileNames, index, index + 1)
						_ = slices.Delete(xdgList, index, index + 1)
					}
				}
			}

			// Sort alphabetically (case insensitive)
			sort.Slice(files, func(x, y int) bool {
				return strings.ToLower(files[x].AppName) < strings.ToLower(files[y].AppName)
			})

			List = append(List, category{
				File: catFile.Name(),
				Name: catNameMatchReference[1],
				Applications: files,
			})
		}
	}

	// Populate the uncategorized list using the items remaining in the xdgList
	for _, item := range xdgList {
		if item.FileName != "" {
			List[0].Applications = append(List[0].Applications, item)
		}
	}

	// Sort alphabetically (case insensitive)
	sort.Slice(List, func(x, y int) bool {
		return strings.ToLower(List[x].File) < strings.ToLower(List[y].File)
	})

	UpdateFileNames()
	UpdateCatNames()
}

func init() {
	Populate()
}
