package attr

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/options"
)

// Support for filtering by some of the attributes.
// Real Windows support is not planned.
const AttrNone FileAttributes = 0

const (
	AttrArchive FileAttributes = 1 << iota
	AttrCompressed
	AttrDevice
	AttrDirectory
	AttrEncrypted
	AttrHidden
	AttrIntegrityStream
	AttrNormal
	AttrNoScrubData
	AttrNotContentIndexed
	AttrOffline
	AttrReadOnly
	AttrReparsePoint
	AttrSparseFile
	AttrSystem
	AttrTemporary
)

type FileAttributes uint32

func (a FileAttributes) Has(attr FileAttributes) bool {
	return a&attr != 0
}

func (a FileAttributes) With(attr FileAttributes) FileAttributes {
	return a | attr
}

func (a FileAttributes) Without(attr FileAttributes) FileAttributes {
	return a &^ attr
}

func (a FileAttributes) String() string {
	var out bytes.Buffer

	if a.Has(AttrDirectory) {
		out.WriteByte('d')
	} else {
		out.WriteByte('-')
	}

	if a.Has(AttrArchive) {
		out.WriteByte('a')
	} else {
		out.WriteByte('-')
	}

	if a.Has(AttrReadOnly) {
		out.WriteByte('r')
	} else {
		out.WriteByte('-')
	}

	if a.Has(AttrHidden) {
		out.WriteByte('h')
	} else {
		out.WriteByte('-')
	}

	if a.Has(AttrSystem) {
		out.WriteByte('s')
	} else {
		out.WriteByte('-')
	}

	if a.Has(AttrReparsePoint) {
		out.WriteByte('l')
	} else {
		out.WriteByte('-')
	}

	return out.String()
}

func ValidateAttributes(attributes []string) error {
	var invalidAttributes []string

	for _, attr := range attributes {
		if attr == "" {
			continue
		}
		if !slices.Contains(options.ValidAttributes, attr) {
			invalidAttributes = append(invalidAttributes, attr)
		}
	}

	if len(invalidAttributes) > 0 {
		return fmt.Errorf(
			"invalid attributes: %s",
			strings.Join(invalidAttributes, ","),
		)
	}

	return nil
}
