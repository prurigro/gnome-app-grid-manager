# Gnome Application Category Manager

An application grid organizer for Gnome.

## Category Folder

A **category folder** is this program's name for the folders of grouped applications in the overview.

## Menu Options

* **Manage application categories**: Interactively move applications between category folders (includes a default **Uncategorized** list)
* **Create new category folder**: Creates a new category folder with no applications inside
* **Delete existing category folder**: Deletes a category folder and moves its applications to **Uncategorized**
* **Clean and sort category files**: Removes applications that aren't installed and visible from each category file, and sorts them alphabetically
* **Apply category folders in Gnome**: Applies the configured category folders and applications to the Gnome overview application grid
* **Restore default layout in Gnome**: Removes category folders from the Gnome overview application grid and resets the layout

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
