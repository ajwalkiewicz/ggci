package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/options"
)

type (
	User  string
	Group string
)

type Item struct {
	Name         string
	Path         string
	ParentDir    string
	ModifiedTime time.Time
	Size         int64
	User         User
	Group        Group
	Mode         ItemMode
}

type ItemMode struct {
	UnixMode   fs.FileMode
	Attributes attr.FileAttributes
}

func NewItem(parentPath string, file fs.FileInfo, options *options.Options) Item {
	var mode fs.FileMode
	var user User
	var group Group
	var path string

	// Technically this is not necessary, because declared variables get
	// their "zeroed" value by default if not defined, but conceptually it
	// is easier to have "positive" if-test then negative
	if options.Name || options.Legacy {
		mode = 0
		user, group = User(""), Group("")
		path = ""
	} else {
		mode = file.Mode()
		user, group = ownerNames(file)
		path = filepath.Join(parentPath, file.Name())
	}

	return Item{
		Name:         displayName(path, file),
		Path:         path,
		ParentDir:    parentPath,
		ModifiedTime: file.ModTime(),
		Size:         file.Size(),
		User:         user,
		Group:        group,
		Mode: ItemMode{
			UnixMode:   mode,
			Attributes: GetFileAttributes(file),
		},
	}
}

func displayName(path string, file fs.FileInfo) string {
	if file.Mode()&fs.ModeSymlink == 0 {
		return file.Name()
	}

	target, err := os.Readlink(path)
	if err != nil {
		return file.Name()
	}

	return fmt.Sprintf("%s -> %s", file.Name(), target)
}

func ownerNames(file fs.FileInfo) (User, Group) {
	stat, ok := file.Sys().(*syscall.Stat_t)
	if !ok {
		return "", ""
	}

	uid := fmt.Sprintf("%d", stat.Uid)
	gid := fmt.Sprintf("%d", stat.Gid)

	userName := uid
	if userInfo, err := user.LookupId(uid); err == nil {
		userName = userInfo.Username
	}

	groupName := gid
	if groupInfo, err := user.LookupGroupId(gid); err == nil {
		groupName = groupInfo.Name
	}

	return User(userName), Group(groupName)
}
