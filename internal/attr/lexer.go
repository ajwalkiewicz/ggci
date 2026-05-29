package attr

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token
	var err error

	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(TokenAnd, l.ch)
	case ',':
		tok = newToken(TokenOr, l.ch)
	case '!':
		tok = newToken(TokenNot, l.ch)
	case 0:
		tok.Type = TokenEOF
		tok.Literal = "EOF"
	default:
		tok.Literal = l.readIdentifier()
		tok.Type, err = LookupKeyword(tok.Literal)
		if err != nil {
			panic(err)
		}

		return tok
	}

	l.readChar()
	return tok
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
