package output

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/filesystem"
)

func TestTableFormatterFormatWritesTable(t *testing.T) {
	listing := filesystem.DirectoryListing{
		Path: "/tmp/example",
		Items: []filesystem.Item{
			{
				Name:         "alpha.txt",
				ModifiedTime: time.Date(2026, time.May, 27, 9, 15, 0, 0, time.UTC),
				Size:         42,
			},
		},
	}
	formatter := TableFormatter{
		Terminal: Terminal{Width: 80},
		View: View{Columns: []Column{
			LegacyLengthColumn,
			NameColumn,
		}},
	}

	var output bytes.Buffer
	if err := formatter.Format(listing, &output); err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "\n    Directory: /tmp/example\n\n      Length Name     \n      ------ ----     \n          42 alpha.txt\n"
	if output.String() != want {
		t.Fatalf("unexpected output:\n%s", output.String())
	}
}

func TestTableFormatterFormatDropsColumnsThatDoNotFit(t *testing.T) {
	listing := filesystem.DirectoryListing{
		Path: "/tmp/example",
		Items: []filesystem.Item{
			{Name: "alpha.txt", Size: 42},
		},
	}
	formatter := TableFormatter{
		Terminal: Terminal{Width: 10},
		View: View{Columns: []Column{
			LegacyLengthColumn,
			NameColumn,
		}},
	}

	var output bytes.Buffer
	if err := formatter.Format(listing, &output); err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	if strings.Contains(output.String(), "Length") {
		t.Fatalf("expected Length column to be dropped:\n%s", output.String())
	}
	if !strings.Contains(output.String(), "Name") {
		t.Fatalf("expected Name column to remain:\n%s", output.String())
	}
}

func TestTableFormatterFormatColorsOnlyNames(t *testing.T) {
	listing := filesystem.DirectoryListing{
		Path: "/tmp/example",
		Items: []filesystem.Item{
			{
				Name:         "scripts",
				ModifiedTime: time.Date(2026, time.May, 27, 9, 15, 0, 0, time.UTC),
				Mode: filesystem.ItemMode{
					UnixMode:   fs.ModeDir | 0o755,
					Attributes: attr.AttrDirectory,
				},
			},
		},
	}
	formatter := TableFormatter{
		Terminal: Terminal{Width: 80},
		Color:    true,
		View: View{Columns: []Column{
			LegacyModeColumn,
			NameColumn,
		}},
	}

	var output bytes.Buffer
	if err := formatter.Format(listing, &output); err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	got := output.String()
	if !strings.Contains(got, directoryStyle+"scripts"+ansiReset) {
		t.Fatalf("expected directory name to be colored:\n%s", got)
	}
	if !strings.Contains(got, tableHeaderStyle+"Name"+ansiReset) {
		t.Fatalf("expected table header to be colored:\n%s", got)
	}
	if strings.Contains(got, directoryStyle+"Name"+ansiReset) {
		t.Fatalf("expected table header to use header color:\n%s", got)
	}
	if !strings.Contains(got, tableHeaderStyle+"----"+ansiReset) {
		t.Fatalf("expected separator to be colored:\n%s", got)
	}
}
