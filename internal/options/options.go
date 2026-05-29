package options

import (
	"bytes"
	"fmt"
)

const (
	ARCHIVE    = "Archive"
	COMPRESSED = "Compressed"
	DEVICE     = "Device"
	DIRECTORY  = "Directory"
	HIDDEN     = "Hidden"
	NORMAL     = "Normal"
	READONLY   = "ReadOnly"
	SYSTEM     = "System"
)

var ValidAttributes = []string{
	ARCHIVE,
	COMPRESSED,
	DEVICE,
	DIRECTORY,
	HIDDEN,
	NORMAL,
	READONLY,
	SYSTEM,
}

type Options struct {
	// Get-ChildItem options
	Path          []string
	Attributes    string
	Exclude       []string
	Include       []string
	Filter        string
	Depth         int64
	File          bool
	Directory     bool
	FollowSymlink bool
	Force         bool
	Hidden        bool
	Name          bool
	ReadOnly      bool
	Recurse       bool
	System        bool

	// GGCI specific options
	Legacy bool
}

func (o *Options) String() string {
	var out bytes.Buffer

	out.WriteString("Options{")
	out.WriteString(fmt.Sprintf("Path: %s, ", o.Path))
	out.WriteString(fmt.Sprintf("Attributes: %q, ", o.Attributes))
	out.WriteString(fmt.Sprintf("Exclude: %s, ", o.Exclude))
	out.WriteString(fmt.Sprintf("Include: %s, ", o.Include))
	out.WriteString(fmt.Sprintf("Filter: %q, ", o.Filter))
	out.WriteString(fmt.Sprintf("Depth: %d, ", o.Depth))
	out.WriteString(fmt.Sprintf("File: %v, ", o.File))
	out.WriteString(fmt.Sprintf("Directory: %v, ", o.Directory))
	out.WriteString(fmt.Sprintf("FollowSymlink: %v, ", o.FollowSymlink))
	out.WriteString(fmt.Sprintf("Force: %v, ", o.Force))
	out.WriteString(fmt.Sprintf("Hidden: %v, ", o.Hidden))
	out.WriteString(fmt.Sprintf("Name: %v, ", o.Name))
	out.WriteString(fmt.Sprintf("ReadOnly: %v, ", o.ReadOnly))
	out.WriteString(fmt.Sprintf("Recurse: %v, ", o.Recurse))
	out.WriteString(fmt.Sprintf("System: %v, ", o.System))
	out.WriteString(fmt.Sprintf("Legacy: %v}", o.Legacy))

	return out.String()
}
