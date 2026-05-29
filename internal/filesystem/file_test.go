package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ajwalkiewicz/ggci/internal/attr"
	"github.com/ajwalkiewicz/ggci/internal/options"
)

func TestMatchAttributes(t *testing.T) {
	dir := t.TempDir()

	filePath := filepath.Join(dir, "regular.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		attributes string
		want       bool
	}{
		{
			name:       "regular file is archive and not directory",
			attributes: "Archive+!Directory",
			want:       true,
		},
		{
			name:       "regular file is not hidden",
			attributes: "Hidden",
			want:       false,
		},
		{
			name:       "read-only spelling matches parser token",
			attributes: "ReadOnly",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchAttributes(info, parseAttributes(t, tt.attributes))
			if got != tt.want {
				t.Fatalf("matchAttributes() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestShouldIncludeFileAllowsHiddenWhenAttributesAskForHidden(t *testing.T) {
	dir := t.TempDir()

	filePath := filepath.Join(dir, ".hidden.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if meetsCriteria(info, emptyOptions(), parseAttributes(t, "")) {
		t.Fatal("hidden file should be excluded by default")
	}

	op := emptyOptions()
	op.Attributes = "Hidden"

	if !meetsCriteria(info, op, parseAttributes(t, op.Attributes)) {
		t.Fatal("hidden file should be included when -Attributes Hidden is used")
	}
}

func parseAttributes(t *testing.T, attributes string) *attr.RootNode {
	t.Helper()

	l := attr.NewLexer(attributes)
	p := attr.NewParser(l)
	node := p.ParseRootNode()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	return &node
}

func emptyOptions() *options.Options {
	return &options.Options{}
}
