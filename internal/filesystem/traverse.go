package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/options"
)

type ListingByDir map[string][]Item

func GetListingByDir(path string, op *options.Options) (ListingByDir, error) {
	listingsByDir := make(ListingByDir)

	root, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	postFileInfoValidator := NewPostFileInfoValidator(op)

	walkFn := func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			// For now: ignore inaccessible paths and continue.
			// Later I may want to collect these errors and print them.
			return nil
		}

		// Do not include the root itself as an item.
		if path == root {
			return nil
		}

		// If -Recurse is disabled, skip anything deeper than direct children.
		if !op.Recurse {
			parent := filepath.Dir(path)

			if parent != root {
				if d.IsDir() {
					return fs.SkipDir
				}

				return nil
			}
		}

		// Pre-FileInfo Validation
		if !op.Force && !op.Hidden && !attributesMentionsHidden(op.Attributes) && isHiddenName(d.Name()) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		fileName := d.Name()

		if op.Legacy {
			fileName = strings.ToUpper(fileName)
		}

		if len(op.Filter) > 0 && !matchFilter(op.Filter, fileName) {
			return nil
		}

		if len(op.Include) > 0 && !matchInclude(op.Include, fileName) {
			return nil
		}

		if len(op.Exclude) > 0 && matchExclude(op.Exclude, fileName) {
			return nil
		}

		// Getting FileInfo
		fi, err := d.Info()
		if err != nil {
			// Could not stat this entry. Continue walking.
			return nil
		}

		// Post-FileInfo && Pre-Item validation
		if !postFileInfoValidator.Match(fi) {
			return nil
		}

		parentDir := filepath.Dir(path)
		item := NewItem(parentDir, fi, op)
		listingsByDir[parentDir] = append(listingsByDir[parentDir], item)

		return nil
	}

	err = filepath.WalkDir(root, walkFn)
	if err != nil {
		return nil, err
	}

	return listingsByDir, nil
}

func GetSingleFileStat(path string, op *options.Options) (ListingByDir, error) {
	listingsByDir := make(ListingByDir, 1)

	root, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	postFileInfoValidator := NewPostFileInfoValidator(op)

	// Skip creating Item and adding if it doesn't match criteria
	if !postFileInfoValidator.Match(fi) {
		return listingsByDir, nil
	}

	parentDir := filepath.Dir(path)
	item := NewItem(parentDir, fi, op)
	listingsByDir[root] = append(listingsByDir[root], item)

	return listingsByDir, nil
}

func isHiddenName(name string) bool {
	return strings.HasPrefix(name, ".")
}

func MapToDirectoryListings(grouped ListingByDir) []DirectoryListing {
	listings := make([]DirectoryListing, 0, len(grouped))

	for dir, items := range grouped {
		var dirItems []Item
		var otherItems []Item

		for _, item := range items {
			if item.Mode.Attributes.Has(attr.AttrDirectory) {
				dirItems = append(dirItems, item)
			} else {
				otherItems = append(otherItems, item)
			}
		}

		slices.SortFunc(dirItems, sortItems)
		slices.SortFunc(otherItems, sortItems)

		listings = append(listings, DirectoryListing{
			Path:  dir,
			Items: append(dirItems, otherItems...),
		})
	}

	return listings
}
