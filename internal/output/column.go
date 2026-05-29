package output

import (
	"strconv"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/filesystem"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

type Column struct {
	Header       string
	MinWidth     int
	MaxWidth     int
	Alignment    Alignment
	CanDrop      bool
	DropPriority int
	Value        func(filesystem.Item) string
}

func Pad(value string, width int, alignment Alignment) string {
	if len(value) >= width {
		return value
	}

	padding := strings.Repeat(" ", width-len(value))

	switch alignment {
	case AlignRight:
		return padding + value
	default:
		return value + padding
	}
}

var NameColumn = Column{
	Header:    "Name",
	MinWidth:  4,
	Alignment: AlignLeft,
	CanDrop:   false,
	Value: func(i filesystem.Item) string {
		return i.Name
	},
}

var UnixSizeColumn = Column{
	Header:       "Size",
	MinWidth:     4,
	MaxWidth:     12,
	Alignment:    AlignRight,
	CanDrop:      true,
	DropPriority: 30,
	Value: func(i filesystem.Item) string {
		return strconv.FormatInt(i.Size, 10)
	},
}

var LegacyLengthColumn = Column{
	Header:       "Length",
	MinWidth:     6,
	MaxWidth:     12,
	Alignment:    AlignRight,
	CanDrop:      true,
	DropPriority: 30,
	Value: func(i filesystem.Item) string {
		return strconv.FormatInt(i.Size, 10)
	},
}

var LastWriteTimeColumn = Column{
	Header:       "LastWriteTime",
	MinWidth:     13,
	MaxWidth:     21,
	Alignment:    AlignRight,
	CanDrop:      true,
	DropPriority: 20,
	Value: func(i filesystem.Item) string {
		return i.ModifiedTime.Format("1/2/2006 15:04")
	},
}

var UnixUserColumn = Column{
	Header:       "User",
	MinWidth:     4,
	MaxWidth:     4,
	Alignment:    AlignLeft,
	CanDrop:      true,
	DropPriority: 40,
	Value: func(i filesystem.Item) string {
		return string(i.User)
	},
}

var UnixGroupColumn = Column{
	Header:       "Group",
	MinWidth:     5,
	MaxWidth:     5,
	Alignment:    AlignLeft,
	CanDrop:      true,
	DropPriority: 40,
	Value: func(i filesystem.Item) string {
		return string(i.Group)
	},
}

var UnixModeColumn = Column{
	Header:       "UnixMode",
	MinWidth:     8,
	MaxWidth:     16,
	Alignment:    AlignLeft,
	CanDrop:      true,
	DropPriority: 50,
	Value: func(i filesystem.Item) string {
		return i.Mode.UnixMode.String()
	},
}

var LegacyModeColumn = Column{
	Header:       "Mode",
	MinWidth:     4,
	MaxWidth:     6,
	Alignment:    AlignLeft,
	CanDrop:      true,
	DropPriority: 50,
	Value: func(i filesystem.Item) string {
		return i.Mode.Attributes.String()
	},
}
