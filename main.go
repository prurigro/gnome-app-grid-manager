package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	desktopFiles = []string{}
	noDisplayFiles = []string{}
)

// Takes an xdg data folder and appends /applications without adding two slashes
func xdgToApplications(dir string) string {
	// Add a slash to the end if one doesn't already exist
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}

	// Append the applications folder
	dir = dir + "applications"

	return dir
}

// Returns the existing directories that contain xdg desktop files
func getDesktopDirs() []string {
	var desktopDirs []string

	// Add the user's applications directory if it exists
	homeDesktopDir := xdgToApplications(os.Getenv("XDG_DATA_HOME"))

	if stat, err := os.Stat(homeDesktopDir); err == nil && stat.IsDir() {
		desktopDirs = append(desktopDirs, homeDesktopDir)
	}

	// Add the other applications directories based on XDG_DATA_DIRS if they exist
	for _, dir := range strings.Split(os.Getenv("XDG_DATA_DIRS"), ":") {
		dir = xdgToApplications(dir)

		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			desktopDirs = append(desktopDirs, dir)
		}
	}

	return desktopDirs
}

func isDesktopFileDisplayed(filePath string) (bool) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Unable to open the file " + filePath)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	noDisplayMatchRegex := regexp.MustCompile("(?i)(nodisplay|hidden|notshowin) *= *(true|gnome)")

	for scanner.Scan() {
		if noDisplayMatchRegex.MatchString(scanner.Text()) {
			return false
		}
	}

	return true
}

func populateDesktopFiles(dir string) {
	items, _ := ioutil.ReadDir(dir)

	for _, item := range items {
		if item.IsDir() {
			populateDesktopFiles(dir + "/" + item.Name())
		} else if !slices.Contains(desktopFiles, item.Name()) && !slices.Contains(noDisplayFiles, item.Name()) {
			desktopFileDisplayed := isDesktopFileDisplayed(dir + "/" + item.Name())

			if desktopFileDisplayed {
				desktopFiles = append(desktopFiles, item.Name())
			} else {
				noDisplayFiles = append(noDisplayFiles, item.Name())
			}
		}
	}
}

func main() {
	// Loop through xdg desktop directories in order of priority and populate the desktopFiles array
	for _, dir := range getDesktopDirs() {
		populateDesktopFiles(dir)
	}

	// Sort alphabetically (case insensitive)
	sort.Slice(desktopFiles, func(x, y int) bool {
		return strings.ToLower(desktopFiles[x]) < strings.ToLower(desktopFiles[y])
	})

	// Print out the list of desktop files
	for _, file := range desktopFiles {
		fmt.Println(file)
	}
}
