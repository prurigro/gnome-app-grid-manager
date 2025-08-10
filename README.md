# Gnome Application Grid Manager

Organize your Gnome overview applications by category

## Menu Options

* **Manage application categories**: Interactively move applications between category folders (includes a default **Uncategorized** list)
* **Create new category folder**: Creates a new category folder with no applications inside
* **Delete existing category folder**: Deletes a category folder and moves its applications to **Uncategorized**
* **Apply category folders in Gnome**: Applies the configured category folders and applications to the Gnome overview application grid
* **Restore default layout in Gnome**: Removes category folders from the Gnome overview application grid and resets the layout
* **Clean and sort config files**: Removes applications that aren't installed and visible from each .category file in the config directory, and sorts them alphabetically

## CLI Options

* **-a**|**--apply**: Apply category folders in Gnome
* **-r**|**--restore**: Restore default layout in Gnome
* **-c**|**--clean**: Clean and sort .category files
* **-h**|**--help**: Show the help text

## Filesystem

This program creates the directory `$XDG_DATA_HOME/gnome-shell/categories` and places `.category` files inside.

Each `.category` file represents a category folder, and contains a list of `.desktop` files that get included inside.

The `$XDG_DATA_HOME` directory and directories in `$XDG_DATA_DIRS` are searched for unique `.desktop` files that aren't configured to be hidden (`NoDisplay=true`, `Hidden=true` or `NotShowIn=gnome`).

## CREDITS

Written by Kevin MacMartin:

* [Forgejo](https://git.darkcloud.ca/kevin)
* [GitHub](https://github.com/prurigro)
* [Arch Linux AUR](https://aur.archlinux.org/packages/?SeB=m&K=prurigro)

## LICENSE

Released under the [MIT license](http://opensource.org/licenses/MIT).
