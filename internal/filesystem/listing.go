package filesystem

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/options"
)

// First version will keep everything in the slice, later on, it needs to
// be optimized to actually work in streaming mode, to not store everything.
var ListingResult []DirectoryListing

type DirectoryListing struct {
	Path  string
	Items []Item
}

func NewDirectoryListing(path string, files []fs.FileInfo, options *options.Options) DirectoryListing {
	var dirItems []Item
	var otherItems []Item

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	for _, file := range files {
		if file.IsDir() {
			dirItems = append(dirItems, NewItem(path, file, options))
		} else {
			otherItems = append(otherItems, NewItem(path, file, options))
		}
	}

	slices.SortFunc(dirItems, sortItems)
	slices.SortFunc(otherItems, sortItems)

	return DirectoryListing{
		Path:  absPath,
		Items: append(dirItems, otherItems...),
	}
}

func sortItems(a, b Item) int {
	first := strings.ToUpper(a.Name)
	second := strings.ToUpper(b.Name)

	return strings.Compare(first, second)
}
