# Gnome Application Category Manager

A TUI program for managing overview app grid folders in Gnome.

## Category Folder

A *category folder* is this program's name for the grouped applications in the overview.

## Menu Options

* **Manage Application Categories**: Choose category folders for installed applications that aren't hidden
* **Create New Category Folder**: Create a new category folder to put applications into
* **Delete Existing Category Folder**: Choose and confirm deletion of a category folder
* **Clean Category Files**: Removes applications that aren't installed and visible from each category file, and sorts them alphabetically
* **Apply Gnome Category Folders**: Applies the configured category folders and applications in the Gnome application grid
* **Clear Gnome Category Folders**: Clears the category folders in the Gnome overview application grid

## Filesystem

This program creates the directory `$XDG_DATA_HOME/gnome-shell/categories` and places `.category` files inside.

Each `.category` file represents a category folder, and contains a list of `.desktop` files that get included inside.

To find `.desktop` files, it searches `$XDG_DATA_HOME` and `$XDG_DATA_DIRS` for valid locations and includes all unique `.desktop` files that aren't configured to be hidden.

## CREDITS

Written by Kevin MacMartin:

* [Forgejo](https://git.darkcloud.ca/kevin)
* [GitHub](https://github.com/prurigro)
* [Arch Linux AUR](https://aur.archlinux.org/packages/?SeB=m&K=prurigro)

## LICENSE

Released under the [MIT license](http://opensource.org/licenses/MIT).
