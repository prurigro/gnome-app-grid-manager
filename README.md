# Gnome Application Category Manager

A TUI program for managing overview app grid folders in Gnome.

## Category Folder

A **category folder** is this program's name for the grouped applications in the overview.

## Menu Options

* **Manage Application Categories**: Interactively move applications between category folders (includes a default **Uncategorized** list)
* **Create New Category Folder**: Creates a new category folder with no applications inside
* **Delete Existing Category Folder**: Deletes a category folder and moves its applications to **Uncategorized**
* **Clean and Sort Category Files**: Removes applications that aren't installed and visible from each category file, and sorts them alphabetically
* **Apply Category Folders in Gnome**: Applies the configured category folders and applications to the Gnome overview application grid
* **Restore Default Layout in Gnome**: Removes category folders from the Gnome overview application grid and resets the layout

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
