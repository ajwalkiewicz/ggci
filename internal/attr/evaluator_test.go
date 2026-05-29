package attr

import "testing"

func TestEvalAttributeExpression(t *testing.T) {
	input := "Archive+!Directory,Hidden"

	tests := []struct {
		name       string
		attributes FileAttributes
		want       bool
	}{
		{
			name:       "archive file matches",
			attributes: AttrArchive,
			want:       true,
		},
		{
			name:       "hidden directory matches because hidden is allowed",
			attributes: AttrDirectory.With(AttrHidden),
			want:       true,
		},
		{
			name:       "archive directory does not match",
			attributes: AttrArchive.With(AttrDirectory),
			want:       false,
		},
		{
			name:       "plain directory does not match",
			attributes: AttrDirectory,
			want:       false,
		},
		{
			name:       "plain file does not match",
			attributes: AttrNone,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evalInput(t, input, tt.attributes)
			testBooleanObject(t, got, tt.want)
		})
	}
}

func evalInput(t *testing.T, input string, attributes FileAttributes) Object {
	t.Helper()

	l := NewLexer(input)
	p := NewParser(l)

	root := p.ParseRootNode()
	checkParserErrors(t, p)

	return Eval(&root, attributes)
}

func testBooleanObject(t *testing.T, obj Object, expected bool) {
	t.Helper()

	result, ok := obj.(*Boolean)
	if !ok {
		t.Fatalf("object is not *Boolean, got=%T (%s)", obj, obj.Inspect())
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value, expected=%t, got=%t", expected, result.Value)
	}
}
