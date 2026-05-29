package attr

import "testing"

func TestParseAttributeExpression(t *testing.T) {
	input := "Archive+!Directory,Hidden"

	l := NewLexer(input)
	p := NewParser(l)

	root := p.ParseRootNode()
	checkParserErrors(t, p)

	if len(root.Nodes) != 1 {
		t.Fatalf("root should contain 1 node, got=%d", len(root.Nodes))
	}

	expression, ok := root.Nodes[0].(*ExpressionNode)
	if !ok {
		t.Fatalf("root.Nodes[0] is not *ExpressionNode, got=%T", root.Nodes[0])
	}

	orNode, ok := expression.Node.(*InfixOperatorNode)
	if !ok {
		t.Fatalf("expression.Node is not *InfixOperatorNode, got=%T", expression.Node)
	}
	if orNode.Operator != "," {
		t.Fatalf("orNode.Operator wrong, expected=\",\", got=%q", orNode.Operator)
	}

	andNode, ok := orNode.Left.(*InfixOperatorNode)
	if !ok {
		t.Fatalf("orNode.Left is not *InfixOperatorNode, got=%T", orNode.Left)
	}
	if andNode.Operator != "+" {
		t.Fatalf("andNode.Operator wrong, expected=\"+\", got=%q", andNode.Operator)
	}

	testAttributeLiteral(t, andNode.Left, TokenArchive)

	notNode, ok := andNode.Right.(*PrefixOperatorNode)
	if !ok {
		t.Fatalf("andNode.Right is not *PrefixOperatorNode, got=%T", andNode.Right)
	}
	if notNode.Operator != "!" {
		t.Fatalf("notNode.Operator wrong, expected=\"!\", got=%q", notNode.Operator)
	}

	testAttributeLiteral(t, notNode.Right, TokenDirectory)
	testAttributeLiteral(t, orNode.Right, TokenHidden)
}

func TestParserCreatesExpectedNodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		nodeType any
	}{
		{
			name:     "attribute literal",
			input:    "Archive",
			nodeType: &AttributeLiteralNode{},
		},
		{
			name:     "prefix operator",
			input:    "!Directory",
			nodeType: &PrefixOperatorNode{},
		},
		{
			name:     "and infix operator",
			input:    "Archive+Hidden",
			nodeType: &InfixOperatorNode{},
		},
		{
			name:     "or infix operator",
			input:    "Archive,Hidden",
			nodeType: &InfixOperatorNode{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := parseSingleNode(t, tt.input)

			switch tt.nodeType.(type) {
			case *AttributeLiteralNode:
				if _, ok := node.(*AttributeLiteralNode); !ok {
					t.Fatalf("node is not *AttributeLiteralNode, got=%T", node)
				}
			case *PrefixOperatorNode:
				if _, ok := node.(*PrefixOperatorNode); !ok {
					t.Fatalf("node is not *PrefixOperatorNode, got=%T", node)
				}
			case *InfixOperatorNode:
				if _, ok := node.(*InfixOperatorNode); !ok {
					t.Fatalf("node is not *InfixOperatorNode, got=%T", node)
				}
			default:
				t.Fatalf("unhandled expected node type %T", tt.nodeType)
			}
		})
	}
}

func parseSingleNode(t *testing.T, input string) Node {
	t.Helper()

	l := NewLexer(input)
	p := NewParser(l)

	root := p.ParseRootNode()
	checkParserErrors(t, p)

	if len(root.Nodes) != 1 {
		t.Fatalf("root should contain 1 node, got=%d", len(root.Nodes))
	}

	expression, ok := root.Nodes[0].(*ExpressionNode)
	if !ok {
		t.Fatalf("root.Nodes[0] is not *ExpressionNode, got=%T", root.Nodes[0])
	}

	return expression.Node
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Fatalf("parser has %d errors: %v", len(errors), errors)
}

func testAttributeLiteral(t *testing.T, node Node, expected TokenType) {
	t.Helper()

	attr, ok := node.(*AttributeLiteralNode)
	if !ok {
		t.Fatalf("node is not *AttributeLiteralNode, got=%T", node)
	}
	if attr.Value != expected {
		t.Fatalf("attr.Value wrong, expected=%q, got=%q", expected, attr.Value)
	}
	if attr.TokenLiteral() != string(expected) {
		t.Fatalf("attr.TokenLiteral wrong, expected=%q, got=%q", expected, attr.TokenLiteral())
	}
}
