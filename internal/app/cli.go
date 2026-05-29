package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/filesystem"
	"github.com/ajwalkiewicz/ggci/internal/options"
	"github.com/ajwalkiewicz/ggci/internal/output"
)

const (
	ExitSuccess int = iota
	ExitGeneralError
	ExitParsingError
	ExitReadingError
)

var version = "dev"

func Run(args []string) int {
	options, err := ParseArgs(args)
	if err != nil {
		fmt.Println(err)
		return ExitParsingError
	}

	// fmt.Println(options)

	// First version will keep everything in the map and slice, later on, it needs to
	// be optimized to actually work in streaming mode, to not store everything.
	var listingsByDir filesystem.ListingByDir

	for _, path := range options.Path {
		info, err := filesystem.GetStats(path, options)
		if err != nil {
			fmt.Printf("Error when reading %s: %s", path, err)
			continue
		}

		if info.IsDir() {
			listingsByDir, err = filesystem.GetListingByDir(path, options)
		} else {
			listingsByDir, err = filesystem.GetSingleFileStat(path, options)
		}

		if err != nil {
			fmt.Printf("Error when reading path %s: %s", path, err)
		}

	}

	var formatter output.Formatter

	formatter = &output.UnixTableFormatter
	if options.Legacy {
		formatter = &output.LegacyTableFormatter
	}

	if options.Name {
		formatter = &output.NameFormatter{}
	}

	listings := filesystem.MapToDirectoryListings(listingsByDir)

	for _, listing := range listings {
		if len(listing.Items) == 0 {
			continue
		}

		if err := formatter.Format(listing, os.Stdout); err != nil {
			fmt.Println("Error writing output:", err)
			return ExitGeneralError
		}
	}

	return ExitSuccess
}

func ParseArgs(args []string) (*options.Options, error) {
	var op options.Options
	var err error

	fs := flag.NewFlagSet("ggci", flag.ExitOnError)

	// Flags that needs to be converted to []string
	tmpExclude := fs.String("Exclude", "", "Exclude entires that match pattern")
	tmpInclude := fs.String("Include", "", "Shows items that match the pattern")

	// String flags
	fs.StringVar(&op.Filter, "Filter", "", "Exclude entries that match the pattern")
	fs.StringVar(
		&op.Attributes,
		"Attributes",
		"",
		fmt.Sprintf("Attributes to filter: %s", strings.Join(options.ValidAttributes, "|")),
	)

	// Int flags
	fs.Int64Var(&op.Depth, "Depth", 0, "Depth of recursion")

	// Boolean Flags
	fs.BoolVar(&op.File, "File", false, "Get a list of files")
	fs.BoolVar(&op.Directory, "Directory", false, "Get a list of files")
	fs.BoolVar(&op.FollowSymlink, "FollowSymlink", false, "Follow symbolic links when using recursion")
	fs.BoolVar(&op.Force, "Force", false, "Show items that are hidden or system files")
	fs.BoolVar(&op.Hidden, "Hidden", false, "Show only hidden files")
	fs.BoolVar(&op.Name, "Name", false, "Get only the names of items")
	fs.BoolVar(&op.ReadOnly, "ReadOnly", false, "Show only read-only items")
	fs.BoolVar(&op.Recurse, "Recurse", false, "Show all items recursively")
	fs.BoolVar(&op.System, "System", false, "Show only system items")
	fs.BoolVar(&op.Legacy, "Legacy", false, "Print output in Windows style")

	showVersion := fs.Bool("Version", false, "print version and exit")

	// Help
	help := fs.Bool("help", false, "print help message")
	fs.Usage = CustomUsage

	fs.Parse(args)

	if *help {
		CustomHelpMessage(fs)
		os.Exit(ExitSuccess)
	}

	if *showVersion {
		fmt.Println("ggci", version)
		os.Exit(0)
	}

	op.Path = fs.Args()

	if *tmpExclude != "" {
		op.Exclude = strings.Split(*tmpExclude, ",")
	}
	if *tmpInclude != "" {
		op.Include = strings.Split(*tmpInclude, ",")
	}

	if op.Legacy {
		op.Filter = strings.ToUpper(op.Filter)

		var newExclude []string
		for _, rule := range op.Exclude {
			newExclude = append(newExclude, strings.ToUpper(rule))
		}
		op.Exclude = newExclude

		var newInclude []string
		for _, rule := range op.Include {
			newInclude = append(newInclude, strings.ToUpper(rule))
		}
		op.Include = newInclude
	}

	// Further custom validations:
	if len(op.Path) == 0 {
		op.Path = []string{"."}
	}

	return &op, err
}

func CustomUsage() {
	fmt.Fprintf(os.Stderr, "Usage: ggci [options] [path]\n")
	fmt.Fprintf(os.Stderr, "  Use --help flag for more information.\n")
}

func CustomHelpMessage(fs *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "ggci - Go Get-ChildItem\n\n")

	fmt.Fprintf(os.Stderr, "Description:\n")
	fmt.Fprintf(os.Stderr, "  Lists files and directories, inspired by PowerShell Get-ChildItem.\n\n")

	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  ggci [options] [path]\n\n")

	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  ggci\n")
	fmt.Fprintf(os.Stderr, "  ggci .\n")
	fmt.Fprintf(os.Stderr, "  ggci -Recurse -Depth 3 /var/log\n\n")

	fmt.Fprintf(os.Stderr, "Options:\n")
	fs.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\nAuthor:\n")
	fmt.Fprintf(os.Stderr, "  Adam Walkiewicz\n")
}
