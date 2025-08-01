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

type Launcher struct {
	AppName string
	FileName string
}

var (
	List = []Launcher{}
	FileNames = []string{}
	AppNames = []string{}
	FileMatchRegex = regexp.MustCompile("(?i) *(.*)\\.desktop *$")
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
func getDesktopFileDirectories() []string {
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
		display bool = true
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
		if name == "" && nameMatchRegex.MatchString(scanner.Text()) {
			name = nameMatchRegex.FindStringSubmatch(scanner.Text())[1]
		}

		if noDisplayFilesMatchRegex.MatchString(scanner.Text()) {
			display = false
		}
	}

	if (name == "") {
		name = filename
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
				List = append(List, Launcher{ AppName: appName, FileName: item.Name() })
				displayFiles = append(displayFiles, item.Name())
			} else {
				noDisplayFiles = append(noDisplayFiles, item.Name())
			}
		}
	}
}

// Returns an array of file names from a list of launchers
func GetFileNames(launchers []Launcher) []string {
	var fileNames []string

	for _, item := range launchers {
		fileNames = append(fileNames, item.FileName)
	}

	return fileNames
}

// Returns an array of app names from a list of launchers
func GetAppNames(launchers []Launcher) []string {
	var appNames []string

	for _, item := range launchers {
		appNames = append(appNames, item.AppName)
	}

	return appNames
}

// Updates the list of all file names
func UpdateFileNames() {
	FileNames = GetFileNames(List)
}

// Updates the list of all app names
func UpdateAppNames() {
	AppNames = GetAppNames(List)
}

// Populate the List of xdg launchers
func Populate() {
	List = nil

	// Loop through xdg desktop directories in order of priority and populate the List array
	for _, dir := range getDesktopFileDirectories() {
		addDirectoryDesktopFiles(dir)
	}

	// Sort alphabetically (case insensitive)
	sort.Slice(List, func(x, y int) bool {
		return strings.ToLower(List[x].AppName) < strings.ToLower(List[y].AppName)
	})

	UpdateFileNames()
	UpdateAppNames()
}

func init() {
	Populate()
}
