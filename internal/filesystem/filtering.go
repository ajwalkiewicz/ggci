package filesystem

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/options"
)

type PostFileInfoValidator struct {
	Options  *options.Options
	FileInfo fs.FileInfo
	RootNode *attr.RootNode
}

func NewPostFileInfoValidator(op *options.Options) *PostFileInfoValidator {
	l := attr.NewLexer(op.Attributes)
	p := attr.NewParser(l)
	node := p.ParseRootNode()

	return &PostFileInfoValidator{
		Options:  op,
		RootNode: &node,
	}
}

func (v *PostFileInfoValidator) Match(fi fs.FileInfo) bool {
	return meetsCriteria(fi, v.Options, v.RootNode)
}

// Check if the FileInfo meets criteria, defined by Attributes and
// flags, except Exclude, Include and Filter
func meetsCriteria(fi fs.FileInfo, op *options.Options, rn *attr.RootNode) bool {
	if len(op.Attributes) > 0 && !matchAttributes(fi, rn) {
		return false
	}

	if op.File && IsFileDirectory(fi) {
		return false
	}

	if op.Directory && !IsFileDirectory(fi) {
		return false
	}

	if op.Hidden && !IsFileHidden(fi) {
		return false
	}

	if !op.Hidden && !op.Force && !attributesMentionsHidden(op.Attributes) && IsFileHidden(fi) {
		return false
	}

	if op.ReadOnly && !IsFileReadOnly(fi) {
		return false
	}

	if op.System && !IsFileSystem(fi) {
		return false
	}

	return true
}

func matchAttributes(fi fs.FileInfo, rn *attr.RootNode) bool {
	fa := GetFileAttributes(fi)
	result := attr.Eval(rn, fa)

	return result.(*attr.Boolean).Value
}

func attributesMentionsHidden(attributes string) bool {
	return strings.Contains(attributes, "Hidden")
}

type PreFileInfoValidator struct {
	Options *options.Options
}

func (v PreFileInfoValidator) Match(name string) bool {
	if len(v.Options.Filter) > 0 && !matchFilter(v.Options.Filter, name) {
		return false
	}

	if len(v.Options.Include) > 0 && !matchInclude(v.Options.Include, name) {
		return false
	}

	if len(v.Options.Exclude) > 0 && matchExclude(v.Options.Exclude, name) {
		return false
	}

	return true
}

func matchExclude(rules []string, name string) bool {
	for _, rule := range rules {
		if match, _ := filepath.Match(rule, name); match {
			return match
		}
	}
	return false
}

func matchInclude(rules []string, name string) bool {
	for _, rule := range rules {
		if match, _ := filepath.Match(rule, name); match {
			return match
		}
	}
	return false
}

func matchFilter(filter string, name string) bool {
	match, _ := filepath.Match(filter, name)
	return match
}
