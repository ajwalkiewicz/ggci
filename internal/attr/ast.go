package attr

import "bytes"

type Node interface {
	TokenLiteral() string
	String() string
}

type RootNode struct {
	Nodes []Node
}

func (r *RootNode) TokenLiteral() string {
	if len(r.Nodes) > 0 {
		return r.Nodes[0].TokenLiteral()
	} else {
		return ""
	}
}

func (r *RootNode) String() string {
	var out bytes.Buffer

	for _, s := range r.Nodes {
		out.WriteString(s.String())
	}

	return out.String()
}

type ExpressionNode struct {
	Token Token
	Node  Node
}

func (en *ExpressionNode) TokenLiteral() string { return en.Token.Literal }

func (en *ExpressionNode) String() string {
	if en.Node != nil {
		return en.Node.String()
	}

	return ""
}

type AttributeLiteralNode struct {
	Token Token
	Value TokenType
}

func (a *AttributeLiteralNode) TokenLiteral() string {
	return a.Token.Literal
}

func (a *AttributeLiteralNode) String() string {
	return a.Token.Literal
}

type PrefixOperatorNode struct {
	Token    Token
	Operator string
	Right    Node
}

func (po *PrefixOperatorNode) TokenLiteral() string { return po.Token.Literal }
func (po *PrefixOperatorNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(po.Operator)
	out.WriteString(po.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixOperatorNode struct {
	Token    Token
	Left     Node
	Operator string
	Right    Node
}

func (io *InfixOperatorNode) TokenLiteral() string { return io.Token.Literal }
func (io *InfixOperatorNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(io.Left.String())
	out.WriteString(" ")
	out.WriteString(io.Operator)
	out.WriteString(" ")
	out.WriteString(io.Right.String())
	out.WriteString(")")

	return out.String()
}
