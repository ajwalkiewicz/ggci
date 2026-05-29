package attr

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	Archive+!Directory,Hidden
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenArchive, "Archive"},
		{TokenAnd, "+"},
		{TokenNot, "!"},
		{TokenDirectory, "Directory"},
		{TokenOr, ","},
		{TokenHidden, "Hidden"},
		{TokenEOF, "EOF"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
