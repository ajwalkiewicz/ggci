package output

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/filesystem"
)

const (
	ansiReset         = "\x1b[0m"
	tableHeaderStyle  = "\x1b[32;1m"
	directoryStyle    = "\x1b[44;1m"
	symbolicLinkStyle = "\x1b[36;1m"
	executableStyle   = "\x1b[32;1m"
	archiveStyle      = "\x1b[31;1m"
	powerShellStyle   = "\x1b[33;1m"
)

var extensionStyles = map[string]string{
	".7z":     archiveStyle,
	".cab":    archiveStyle,
	".gz":     archiveStyle,
	".nupkg":  archiveStyle,
	".tar":    archiveStyle,
	".tgz":    archiveStyle,
	".zip":    archiveStyle,
	".ps1":    powerShellStyle,
	".ps1xml": powerShellStyle,
	".psd1":   powerShellStyle,
	".psm1":   powerShellStyle,
}

func itemStyle(item filesystem.Item) string {
	attributes := item.Mode.Attributes

	if attributes.Has(attr.AttrReparsePoint) {
		return symbolicLinkStyle
	}
	if attributes.Has(attr.AttrDirectory) {
		return directoryStyle
	}
	if item.Mode.UnixMode.IsRegular() && item.Mode.UnixMode.Perm()&0o111 != 0 {
		return executableStyle
	}
	if style, ok := extensionStyles[strings.ToLower(filepath.Ext(item.Name))]; ok {
		return style
	}

	return ""
}

func decoratePaddedValue(raw string, padded string, style string) string {
	if raw == "" || style == "" {
		return padded
	}

	rawIndex := strings.Index(padded, raw)
	if rawIndex < 0 {
		return padded
	}

	rawEnd := rawIndex + len(raw)
	return padded[:rawIndex] + style + raw + ansiReset + padded[rawEnd:]
}

func ShouldColor(file *os.File) bool {
	if file == nil {
		return false
	}
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}

	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&fs.ModeCharDevice != 0
}
