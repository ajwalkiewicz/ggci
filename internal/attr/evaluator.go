package attr

import "fmt"

type ObjectType string

const (
	BOOLEAN_OBJ = "BOOLEAN"
	ERROR_OBJ   = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func Eval(node Node, fileAttributes FileAttributes) Object {
	switch node := node.(type) {
	case *RootNode:
		return evalRootNode(node, fileAttributes)
	case *ExpressionNode:
		return Eval(node.Node, fileAttributes)
	case *AttributeLiteralNode:
		return evalAttributeLiteral(node, fileAttributes)
	case *PrefixOperatorNode:
		right := Eval(node.Right, fileAttributes)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixOperatorNode:
		left := Eval(node.Left, fileAttributes)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, fileAttributes)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func evalRootNode(node *RootNode, fileAttributes FileAttributes) Object {
	var result Object

	for _, node := range node.Nodes {
		result = Eval(node, fileAttributes)

		switch result := result.(type) {
		case *Error:
			return result
		}
	}

	return result
}

func evalAttributeLiteral(node Node, fa FileAttributes) Object {
	switch node.TokenLiteral() {
	case TokenArchive:
		return nativeBoolToBooleanObject(fa.Has(AttrArchive))
	case TokenCompressed:
		return nativeBoolToBooleanObject(fa.Has(AttrCompressed))
	case TokenDevice:
		return nativeBoolToBooleanObject(fa.Has(AttrDevice))
	case TokenDirectory:
		return nativeBoolToBooleanObject(fa.Has(AttrDirectory))
	case TokenEncrypted:
		return nativeBoolToBooleanObject(fa.Has(AttrEncrypted))
	case TokenHidden:
		return nativeBoolToBooleanObject(fa.Has(AttrHidden))
	case TokenIntegrityStream:
		return nativeBoolToBooleanObject(fa.Has(AttrIntegrityStream))
	case TokenNormal:
		return nativeBoolToBooleanObject(fa.Has(AttrNormal))
	case TokenNoScrubData:
		return nativeBoolToBooleanObject(fa.Has(AttrNoScrubData))
	case TokenNotContentIndexed:
		return nativeBoolToBooleanObject(fa.Has(AttrNotContentIndexed))
	case TokenOffline:
		return nativeBoolToBooleanObject(fa.Has(AttrOffline))
	case TokenReadOnly:
		return nativeBoolToBooleanObject(fa.Has(AttrReadOnly))
	case TokenReparsePoint:
		return nativeBoolToBooleanObject(fa.Has(AttrReparsePoint))
	case TokenSparseFile:
		return nativeBoolToBooleanObject(fa.Has(AttrSparseFile))
	case TokenSystem:
		return nativeBoolToBooleanObject(fa.Has(AttrSystem))
	case TokenTemporary:
		return nativeBoolToBooleanObject(fa.Has(AttrTemporary))
	default:
		return newError("unknown attribute: %s", node.TokenLiteral())
	}
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalNotOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalNotOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return newError("unknown operator: %s%s", right, right.Type())
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	if left.Type() != BOOLEAN_OBJ && right.Type() != BOOLEAN_OBJ {
		return newError("Only booleans are supported")
	}

	switch {
	case operator == "+":
		return nativeBoolToBooleanObject(left.(*Boolean).Value && right.(*Boolean).Value)
	case operator == ",":
		return nativeBoolToBooleanObject(left.(*Boolean).Value || right.(*Boolean).Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
