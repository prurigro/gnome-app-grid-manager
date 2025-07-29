package main

import (
	"fmt"
	"git.darkcloud.ca/kevin/gnome-appcat-manager/xdg"
)

func main() {
	// Print out the list of desktop files
	for _, file := range xdg.DisplayFiles {
		fmt.Println(file)
	}
}
