package categories

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	// "slices"
	// "sort"
	// "strings"
)

type category struct {
	file string
	name string
	applications []string
}

var (
	Categories = []category{}
	categoriesDirectory = os.Getenv("XDG_DATA_HOME") + "/gnome-shell/categories"
)

func createCatDirWhenMissing() {
	if stat, err := os.Stat(categoriesDirectory); err != nil || !stat.IsDir() {
		os.MkdirAll(categoriesDirectory, 0755)
	}
}

// Returns an array of file names
func FileNames() []string {
	var fileNames []string

	for _, item := range Categories {
		fileNames = append(fileNames, item.file)
	}

	return fileNames
}

// Returns an array of app names
func CatNames() []string {
	var catNames []string

	for _, item := range Categories {
		catNames = append(catNames, item.name)
	}

	return catNames
}

func Populate() {
	var files []string

	createCatDirWhenMissing()
	categoryFiles, _ := os.ReadDir(categoriesDirectory)

	// Regex to pull just the name out of the category file
	catNameRegex := regexp.MustCompile("(?i)^ *(.*)\\.category *$")

	// Regex to ensure we're only looking at .desktop files
	xdgMatchRegex := regexp.MustCompile("(?i)^ *(.*)\\.desktop *$")

	for _, catFile := range categoryFiles {
		if catNameRegex.MatchString(catFile.Name()) {
			files = nil
			catNameMatchReference := catNameRegex.FindStringSubmatch(catFile.Name())
			filePath := categoriesDirectory + "/" + catFile.Name()
			file, err := os.Open(filePath)

			if err != nil {
				log.Fatal("Unable to open the file " + filePath)
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if xdgMatchRegex.MatchString(scanner.Text()) {
					files = append(files, xdgMatchRegex.FindStringSubmatch(scanner.Text())[0])
				}
			}

			Categories = append(Categories, category{ file: catFile.Name(), name: catNameMatchReference[1], applications: []string{} })

			// // Sort alphabetically (case insensitive)
			// sort.Slice(list, func(x, y int) bool {
			// 	return strings.ToLower(list[x]) < strings.ToLower(list[y])
			// })
		}
	}

	fmt.Println(Categories)
}

func init() {
	Populate()
}
