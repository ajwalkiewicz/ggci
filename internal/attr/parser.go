package attr

import "fmt"

const (
	_ int = iota
	LOWEST
	OR
	AND
	NOT
)

var precedences = map[TokenType]int{
	TokenNot: NOT,
	TokenOr:  OR,
	TokenAnd: AND,
}

type (
	prefixParseFn func() Node
	infixParseFn  func(Node) Node
)

type Parser struct {
	l *Lexer

	curToken  Token
	peekToken Token
	errors    []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.infixParseFns = make(map[TokenType]infixParseFn)

	// Attributes
	p.registerPrefix(TokenArchive, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenCompressed, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenDevice, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenDirectory, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenEncrypted, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenHidden, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenIntegrityStream, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenNormal, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenNoScrubData, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenNotContentIndexed, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenOffline, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenReadOnly, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenReparsePoint, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenSparseFile, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenSystem, p.parseAttributeLiteralNode)
	p.registerPrefix(TokenTemporary, p.parseAttributeLiteralNode)

	// Operators
	p.registerPrefix(TokenNot, p.parsePrefixOperatorNode)
	p.registerInfix(TokenOr, p.parseInfixOperatorNode)
	p.registerInfix(TokenAnd, p.parseInfixOperatorNode)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseRootNode() RootNode {
	rootNode := RootNode{}
	rootNode.Nodes = []Node{}

	for !p.curTokenIs(TokenEOF) {
		node := p.parseExpressionNode()
		if node != nil {
			rootNode.Nodes = append(rootNode.Nodes, node)
		}
		p.nextToken()
	}

	return rootNode
}

func (p *Parser) parseExpressionNode() *ExpressionNode {
	node := &ExpressionNode{Token: p.curToken}
	node.Node = p.parseNode(LOWEST)

	return node
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parsePrefixOperatorNode() Node {
	// defer untrace(trace("parsePrefixExpression"))
	expression := &PrefixOperatorNode{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseNode(NOT)

	return expression
}

func (p *Parser) parseInfixOperatorNode(left Node) Node {
	node := &InfixOperatorNode{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	node.Right = p.parseNode(precedence)

	return node
}

func (p *Parser) parseAttributeLiteralNode() Node {
	node := &AttributeLiteralNode{
		Token: p.curToken,
		Value: p.curToken.Type,
	}
	return node
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseNode(precedence int) Node {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftNode := prefix()

	for !p.peekTokenIs(TokenEOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftNode
		}

		p.nextToken()

		leftNode = infix(leftNode)
	}

	return leftNode
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next symbol to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
