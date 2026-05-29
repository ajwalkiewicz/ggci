package attr

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	TokenIllegal = "ILLEGAL"
	TokenEOF     = "EOF"

	// Operators
	TokenNot = "!"
	TokenAnd = "+"
	TokenOr  = ","

	// Attribute Keywords
	TokenArchive           = "Archive"
	TokenCompressed        = "Compressed"
	TokenDevice            = "Device"
	TokenDirectory         = "Directory"
	TokenEncrypted         = "Encrypted"
	TokenHidden            = "Hidden"
	TokenIntegrityStream   = "IntegrityStream"
	TokenNormal            = "Normal"
	TokenNoScrubData       = "NoScrubData"
	TokenNotContentIndexed = "NotContentIndexed"
	TokenOffline           = "Offline"
	TokenReadOnly          = "ReadOnly"
	TokenReparsePoint      = "ReparsePoint"
	TokenSparseFile        = "SparseFile"
	TokenSystem            = "System"
	TokenTemporary         = "Temporary"
)

var attributeKeywords = map[string]TokenType{
	"Archive":           TokenArchive,
	"Compressed":        TokenCompressed,
	"Device":            TokenDevice,
	"Directory":         TokenDirectory,
	"Encrypted":         TokenEncrypted,
	"Hidden":            TokenHidden,
	"IntegrityStream":   TokenIntegrityStream,
	"Normal":            TokenNormal,
	"NoScrubData":       TokenNoScrubData,
	"NotContentIndexed": TokenNotContentIndexed,
	"Offline":           TokenOffline,
	"ReadOnly":          TokenReadOnly,
	"ReparsePoint":      TokenReparsePoint,
	"SparseFile":        TokenSparseFile,
	"System":            TokenSystem,
	"Temporary":         TokenTemporary,
}

func LookupKeyword(attr string) (TokenType, error) {
	if tok, ok := attributeKeywords[attr]; ok {
		return tok, nil
	}

	return TokenType(""), fmt.Errorf("%q is not on the list of valid attribute", attr)
}
