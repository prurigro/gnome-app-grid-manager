package application

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"github.com/prurigro/gnome-app-grid-manager/color"
	"github.com/prurigro/gnome-app-grid-manager/env"
)

type Data struct {
	File string
	Name string
}

var (
	List = []Data{}
	FileMatchRegex = regexp.MustCompile("(?i) *(.*)\\.desktop *$")
	displayFiles = []string{}
	noDisplayFiles = []string{}
)

// Takes an xdg data folder and appends /applications without adding two slashes
func xdgDataToApplications(dir string) string {
	// Add a slash to the end if one doesn't already exist
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}

	// Append the applications folder
	dir = dir + "applications"

	return dir
}

// Returns the existing directories that contain xdg desktop files
func getDesktopFileDirectories() []string {
	var (
		desktopDirs []string
	)

	// Add the user's applications directory if it exists
	homeDesktopDir := xdgDataToApplications(env.XdgDataHome)

	if stat, err := os.Stat(homeDesktopDir); err == nil && stat.IsDir() {
		desktopDirs = append(desktopDirs, homeDesktopDir)
	}

	// Add the other applications directories based on XDG_DATA_DIRS if they exist
	for _, dir := range strings.Split(env.XdgDataDirs, ":") {
		dir = xdgDataToApplications(dir)

		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			desktopDirs = append(desktopDirs, dir)
		}
	}

	return desktopDirs
}

// Checks to see if a .desktop file is configured with NoDisplay=true, Hidden=true or NotShowIn=gnome
func getDesktopFileMeta(dir string, file string) (string, bool) {
	var (
		filePath string = dir + "/" + file
		display bool = true
		name string = ""
	)

	fileData, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Unable to open the file " + color.Red(filePath))
	}

	defer fileData.Close()

	nameMatchRegex := regexp.MustCompile("(?i)^ *name *= *(.*) *$")
	noDisplayFilesMatchRegex := regexp.MustCompile("(?i)^ *(nodisplay|hidden|notshowin) *= *(true|gnome)")
	scanner := bufio.NewScanner(fileData)

	for scanner.Scan() {
		if name == "" && nameMatchRegex.MatchString(scanner.Text()) {
			name = nameMatchRegex.FindStringSubmatch(scanner.Text())[1]
		}

		if noDisplayFilesMatchRegex.MatchString(scanner.Text()) {
			display = false
		}
	}

	if (name == "") {
		name = file
	}

	return name, display
}

// Adds unseen .desktop files in a given directory and below to either displayFiles or noDisplayFiles
func addDirectoryDesktopFiles(dir string) {
	items, _ := os.ReadDir(dir)

	for _, item := range items {
		if item.IsDir() {
			addDirectoryDesktopFiles(dir + "/" + item.Name())
		} else if FileMatchRegex.MatchString(item.Name()) && !slices.Contains(displayFiles, item.Name()) && !slices.Contains(noDisplayFiles, item.Name()) {
			appName, displayed := getDesktopFileMeta(dir, item.Name())

			if displayed {
				List = append(List, Data{ Name: appName, File: item.Name() })
				displayFiles = append(displayFiles, item.Name())
			} else {
				noDisplayFiles = append(noDisplayFiles, item.Name())
			}
		}
	}
}

// Returns an array of file names from a list of applications
func GetFiles(apps []Data) []string {
	var files []string

	for _, item := range apps {
		files = append(files, item.File)
	}

	return files
}

// Returns an array of app names from a list of applications
func GetNames(apps []Data) []string {
	var names []string

	for _, item := range apps {
		names = append(names, item.Name + " " + color.Gray(item.File))
	}

	return names
}

// Sort a list of applications
func Sort(l *[]Data) {
	sort.Slice(*l, func(x, y int) bool {
		return strings.ToLower((*l)[x].Name) < strings.ToLower((*l)[y].Name)
	})
}

// Populate the List of applications
func Populate() {
	// Reset the list of applications
	List = nil

	// Loop through xdg desktop directories in order of priority and populate the List array
	for _, dir := range getDesktopFileDirectories() {
		addDirectoryDesktopFiles(dir)
	}

	// Sort alphabetically (case insensitive)
	Sort(&List)
}

func init() {
	Populate()
}
