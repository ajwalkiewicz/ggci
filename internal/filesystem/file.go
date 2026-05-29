package filesystem

import (
	"io/fs"
	"os"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/options"
)

func GetStats(path string, op *options.Options) (fs.FileInfo, error) {
	var info fs.FileInfo
	var err error

	if op.FollowSymlink {
		info, err = os.Stat(path)
	} else {
		info, err = os.Lstat(path)
	}

	if err != nil {
		return nil, err
	}

	return info, nil
}

// For Windows a regular file is an "Archive"
func IsFileArchive(f fs.FileInfo) bool {
	return !f.IsDir()
}

// Not Supported
func IsFileCompressed(f fs.FileInfo) bool {
	return false
}

func IsFileDevice(f fs.FileInfo) bool {
	return f.Mode()&fs.ModeDevice != 0
}

func IsFileDirectory(f fs.FileInfo) bool {
	return f.IsDir()
}

// Not Supported
func IsFileEncrypted(f fs.FileInfo) bool {
	return false
}

// We can emulate hidden attribute by checking if the file name starts
// with dot '.'
func IsFileHidden(f fs.FileInfo) bool {
	return strings.HasPrefix(f.Name(), ".")
}

// Not Supported
func IsFileIntegrityStream(f fs.FileInfo) bool { return false }

// TODO: Questionable
func IsFileNormal(f fs.FileInfo) bool {
	return IsFileArchive(f) && !IsFileDevice(f)
}

// Not Supported
func IsFileNoScrubData(f fs.FileInfo) bool { return false }

// Not Supported
func IsFileNotContentIndexed(f fs.FileInfo) bool { return false }

// Not Supported
func IsFileOffline(f fs.FileInfo) bool { return false }

// We can emulate ReadOnly attribute by checking if file has only
// read permission in all groups
func IsFileReadOnly(f fs.FileInfo) bool {
	return f.Mode().Perm()&0o222 == 0
}

// 'ReparsePoint' is more or less equivalent of symlink on Linux
func IsFileReparsePoint(f fs.FileInfo) bool {
	return f.Mode()&fs.ModeSymlink != 0
}

// Not Supported
func IsFileSparseFile(f fs.FileInfo) bool { return false }

// Not Supported
func IsFileTemporary(f fs.FileInfo) bool { return false }

// Not Supported
func IsFileSystem(f fs.FileInfo) bool { return false }

// Not supported attributes are not included in the function.
// So effectively this checks following attributes:
// - Archive
// - Device
// - Directory
// - Hidden
// - ReadOnly
// - ReparsePoint
func GetFileAttributes(file fs.FileInfo) attr.FileAttributes {
	var attributes attr.FileAttributes

	if IsFileArchive(file) {
		attributes = attributes.With(attr.AttrArchive)
	}

	if IsFileDevice(file) {
		attributes = attributes.With(attr.AttrDevice)
	}

	if IsFileDirectory(file) {
		attributes = attributes.With(attr.AttrDirectory)
	}

	if IsFileHidden(file) {
		attributes = attributes.With(attr.AttrHidden)
	}

	if IsFileReadOnly(file) {
		attributes = attributes.With(attr.AttrReadOnly)
	}

	if IsFileReparsePoint(file) {
		attributes = attributes.With(attr.AttrReparsePoint)
	}

	return attributes
}
