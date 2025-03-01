// Package eval handles the evaluation of Lox source code and implements the
// tree walking interpreter.
package eval

import (
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/glox/internal/syntax/types"
)

// Eval evaluates a Lox AST Node.
func Eval(node ast.Node) (types.Type, error) {
	switch node := node.(type) {
	case ast.Program:
		return evalStatements(node.Statements)
	case ast.ExpressionStatement:
		return Eval(node.Value)
	case ast.GroupedExpression:
		return Eval(node.Value)
	case ast.UnaryExpression:
		return evalUnaryExpression(node)
	case ast.Number:
		return evalNumber(node), nil
	case ast.Bool:
		return evalBool(node), nil
	default:
		return nil, fmt.Errorf("unhandled ast.Node in Eval: %T", node)
	}
}

func evalStatements(statements []ast.Statement) (types.Type, error) {
	var result types.Type
	var err error
	for _, statement := range statements {
		result, err = Eval(statement)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func evalBool(node ast.Bool) types.Type {
	if node.Value {
		return types.True
	}
	return types.False
}

func evalNumber(node ast.Number) types.Type {
	return types.Number{Value: node.Value}
}

func evalUnaryExpression(node ast.UnaryExpression) (types.Type, error) {
	operand, err := Eval(node.Value)
	if err != nil {
		return nil, err
	}

	switch node.Tok.Kind {
	case token.Minus:
		if operand.Kind() != types.KindNumber {
			return nil, fmt.Errorf("cannot negate a non numeric type: %s", operand)
		}
		number, ok := operand.(types.Number)
		if !ok {
			return nil, fmt.Errorf("could not cast %T to Number", operand)
		}
		return types.Number{Value: -number.Value}, nil
	case token.Bang:
		return types.Bool{Value: !types.IsTruthy(operand)}, nil
	default:
		return nil, fmt.Errorf("unsupported unary operator: %s", node.Tok.Kind.Lexeme())
	}
}
