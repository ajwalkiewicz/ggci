package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/filesystem"
)

type Formatter interface {
	Format(dl filesystem.DirectoryListing, w io.Writer) error
}

type TableFormatter struct {
	Terminal Terminal
	View     View
	Color    bool
}

func (tf *TableFormatter) Format(dl filesystem.DirectoryListing, w io.Writer) error {
	columns, widths := tf.fitColumns(dl)

	if _, err := fmt.Fprintf(w, "\n    Directory: %s\n\n", dl.Path); err != nil {
		return err
	}

	if len(columns) == 0 {
		return nil
	}

	if err := writeHeaderRow(w, columns, widths, tf.Color); err != nil {
		return err
	}

	if err := writeSeparator(w, columns, widths, tf.Color); err != nil {
		return err
	}

	for _, item := range dl.Items {
		if err := writeItemRow(w, columns, widths, item, tf.Color); err != nil {
			return err
		}
	}

	return nil
}

func (tf *TableFormatter) fitColumns(dl filesystem.DirectoryListing) ([]Column, []int) {
	columns := append([]Column(nil), tf.View.Columns...)
	widths := columnWidths(dl, columns)
	if len(columns) == 0 {
		return columns, widths
	}

	width := tf.Terminal.Width
	if width <= 0 {
		return columns, widths
	}

	for tableWidth(widths) > width {
		shrinkIndex := nextShrinkIndex(columns, widths)
		if shrinkIndex >= 0 {
			widths[shrinkIndex]--
			continue
		}

		dropIndex := nextDropIndex(columns)
		if dropIndex < 0 {
			break
		}
		columns = append(columns[:dropIndex], columns[dropIndex+1:]...)
		widths = append(widths[:dropIndex], widths[dropIndex+1:]...)
	}

	return columns, widths
}

func nextShrinkIndex(columns []Column, widths []int) int {
	for i := len(columns) - 1; i >= 0; i-- {
		if widths[i] > columns[i].MinWidth {
			return i
		}
	}

	return -1
}

func nextDropIndex(columns []Column) int {
	for i := len(columns) - 1; i >= 0; i-- {
		column := columns[i]
		if !column.CanDrop {
			continue
		}
		if len(columns) == 1 {
			continue
		}
		return i
	}

	return -1
}

func columnWidths(dl filesystem.DirectoryListing, columns []Column) []int {
	widths := make([]int, len(columns))

	for i, column := range columns {
		width := max(len(column.Header), column.MinWidth)

		if column.MaxWidth > 0 {
			width = max(width, column.MaxWidth)
		} else {
			for _, item := range dl.Items {
				width = max(width, len(column.Value(item)))
			}
		}

		widths[i] = width
	}

	return widths
}

func tableWidth(widths []int) int {
	if len(widths) == 0 {
		return 0
	}

	width := 0
	for _, columnWidth := range widths {
		width += columnWidth
	}

	return width + columnGapWidth(len(widths))
}

func columnGapWidth(columnCount int) int {
	if columnCount <= 1 {
		return 0
	}

	return columnCount - 1
}

func writeRow(w io.Writer, columns []Column, widths []int, value func(Column) string) error {
	return writeDecoratedRow(w, columns, widths, value, nil)
}

func writeHeaderRow(w io.Writer, columns []Column, widths []int, color bool) error {
	return writeDecoratedRow(
		w,
		columns,
		widths,
		func(c Column) string {
			return c.Header
		},
		func(_ Column, raw string, padded string) string {
			if !color {
				return padded
			}

			return decoratePaddedValue(raw, padded, tableHeaderStyle)
		},
	)
}

func writeItemRow(
	w io.Writer,
	columns []Column,
	widths []int,
	item filesystem.Item,
	color bool,
) error {
	return writeDecoratedRow(
		w,
		columns,
		widths,
		func(c Column) string {
			return c.Value(item)
		},
		func(c Column, raw string, padded string) string {
			if !color || c.Header != NameColumn.Header {
				return padded
			}

			return decoratePaddedValue(raw, padded, itemStyle(item))
		},
	)
}

func writeDecoratedRow(
	w io.Writer,
	columns []Column,
	widths []int,
	value func(Column) string,
	decorate func(Column, string, string) string,
) error {
	parts := make([]string, len(columns))

	for i, column := range columns {
		raw := truncate(value(column), widths[i])
		padded := Pad(raw, widths[i], column.Alignment)
		if decorate != nil {
			padded = decorate(column, raw, padded)
		}
		parts[i] = padded
	}

	_, err := fmt.Fprintln(w, strings.Join(parts, " "))
	return err
}

func writeSeparator(w io.Writer, columns []Column, widths []int, color bool) error {
	parts := make([]string, len(columns))

	for i, column := range columns {
		separatorWidth := min(len(column.Header), widths[i])
		raw := strings.Repeat("-", separatorWidth)
		padded := Pad(raw, widths[i], column.Alignment)
		if color {
			padded = decoratePaddedValue(raw, padded, tableHeaderStyle)
		}
		parts[i] = padded
	}

	_, err := fmt.Fprintln(w, strings.Join(parts, " "))
	return err
}

func truncate(value string, width int) string {
	if width <= 0 || len(value) <= width {
		return value
	}

	return value[:width]
}

type NameFormatter struct{}

func (nf *NameFormatter) Format(dl filesystem.DirectoryListing, w io.Writer) error {
	for _, item := range dl.Items {
		if _, err := fmt.Fprintf(w, "%s\r\n", item.Name); err != nil {
			return err
		}
	}

	return nil
}

var UnixTableFormatter = TableFormatter{
	Terminal: NewTerminal(),
	View:     UnixView,
	Color:    ShouldColor(os.Stdout),
}

var LegacyTableFormatter = TableFormatter{
	Terminal: NewTerminal(),
	View:     LegacyView,
	Color:    ShouldColor(os.Stdout),
}
