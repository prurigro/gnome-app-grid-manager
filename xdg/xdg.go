package xdg

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

type xdgMeta struct {
	appName string
	fileName string
}

var (
	FilesMeta = []xdgMeta{}
	displayFiles = []string{}
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

// Checks to see if a .desktop file is configured with NoDisplay=true, Hidden=true or NotShowIn=gnome
func getDesktopFileMeta(dir string, filename string) (string, bool) {
	var (
		filePath string = dir + "/" + filename
		display bool = false
		name string = ""
	)

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Unable to open the file " + filePath)
	}

	defer file.Close()

	nameMatchRegex := regexp.MustCompile("(?i)^ *name *= *(.*) *$")
	noDisplayFilesMatchRegex := regexp.MustCompile("(?i)^ *(nodisplay|hidden|notshowin) *= *(true|gnome)")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if nameMatchRegex.MatchString(scanner.Text()) {
			nameMatchReference := nameMatchRegex.FindStringSubmatch(scanner.Text())
			name = nameMatchReference[1]
		} else if noDisplayFilesMatchRegex.MatchString(scanner.Text()) {
			display = true
		}
	}

	if (name == "") {
		name = filename
	}

	return name, display
}

// Adds unseen .desktop files in a given directory and below to either displayFiles or noDisplayFiles
func populateDesktopFiles(dir string) {
	items, _ := os.ReadDir(dir)

	for _, item := range items {
		if item.IsDir() {
			populateDesktopFiles(dir + "/" + item.Name())
		} else if !slices.Contains(displayFiles, item.Name()) && !slices.Contains(noDisplayFiles, item.Name()) {
			appName, displayed := getDesktopFileMeta(dir, item.Name())

			if displayed {
				FilesMeta = append(FilesMeta, xdgMeta{ appName: appName, fileName: item.Name() })
				displayFiles = append(displayFiles, item.Name())
			} else {
				noDisplayFiles = append(noDisplayFiles, item.Name())
			}
		}
	}
}

// Returns an array of file names
func FileNames() []string {
	var fileNames []string

	for _, file := range FilesMeta {
		fileNames = append(fileNames, file.fileName)
	}

	return fileNames
}

// Returns an array of app names
func AppNames() []string {
	var appNames []string

	for _, file := range FilesMeta {
		appNames = append(appNames, file.appName)
	}

	return appNames
}

func init() {
	// Loop through xdg desktop directories in order of priority and populate the FilesMeta array
	for _, dir := range getDesktopDirs() {
		populateDesktopFiles(dir)
	}

	// Sort alphabetically (case insensitive)
	sort.Slice(FilesMeta, func(x, y int) bool {
		return strings.ToLower(FilesMeta[x].appName) < strings.ToLower(FilesMeta[y].appName)
	})
}
