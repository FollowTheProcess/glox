// Package eval handles the evaluation of Lox source code and implements the
// tree walking interpreter.
package eval

import (
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/glox/internal/syntax/types"
)

// TODO(@FollowTheProcess): Currently it stops on the first error, we should report the error
// to the user with position info and then move on to the next expression. Maybe need an interpreter
// struct with fields tracking errors etc.

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
	case ast.BinaryExpression:
		return evalBinaryExpression(node)
	case ast.Number:
		return evalNumber(node), nil
	case ast.Bool:
		return evalBool(node), nil
	case ast.String:
		return evalString(node), nil
	default:
		return nil, fmt.Errorf("unhandled ast.Node in Eval: %T", node)
	}
}

// evalStatements iterates through all the statements in the program, evaluating each
// and returning the final type.
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

// evalBool interprets a bool AST node, returning either the True or False singleton.
func evalBool(node ast.Bool) *types.Bool {
	if node.Value {
		return types.True
	}
	return types.False
}

// evalNumber interprets a numeric AST node, returning a types.Number.
func evalNumber(node ast.Number) types.Number {
	return types.Number{Value: node.Value}
}

// evalString interprets a string AST node, returning a types.String.
func evalString(node ast.String) types.String {
	return types.String{Value: node.Value}
}

// evalUnaryExpression interprets a unary expression like `-5` or `!true`.
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

// evalBinaryExpression interprets a binary expression like `5 + 5` or `x != y`.
func evalBinaryExpression(node ast.BinaryExpression) (types.Type, error) {
	left, err := Eval(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := Eval(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Op.Kind {
	case token.Minus:
		return evalBinarySubtract(left, right)
	case token.Plus:
		return evalBinaryAdd(left, right)
	case token.ForwardSlash:
		return evalBinaryDivide(left, right)
	case token.Star:
		return evalBinaryMultiply(left, right)
	case token.GreaterThan:
		return evalGreaterThan(left, right)
	case token.GreaterThanEq:
		return evalGreaterThanEq(left, right)
	case token.LessThan:
		return evalLessThan(left, right)
	case token.LessThanEq:
		return evalLessThanEq(left, right)
	default:
		return nil, fmt.Errorf("unsupported binary operator: %s", node.Op.Kind.Lexeme())
	}
}

// evalBinarySubtract interprets a binary subtraction e.g. `5 - 3`.
func evalBinarySubtract(left, right types.Type) (types.Number, error) {
	var zero types.Number
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return zero, err
	}

	result := types.Number{Value: l.Value - r.Value}
	return result, nil
}

// evalBinaryAdd interprets a binary addition e.g. `x + y`.
//
// It is overloaded in the case of two strings to concat the string.
func evalBinaryAdd(left, right types.Type) (types.Type, error) {
	switch left := left.(type) {
	case types.String:
		// Make sure right is also a string
		r, ok := right.(types.String)
		if !ok {
			return nil, fmt.Errorf("invalid types for binary add: left (String) + right (%T)", right)
		}
		return types.String{Value: left.Value + r.Value}, nil
	default:
		l, r, err := checkNumeric(left, right)
		if err != nil {
			return nil, err
		}
		return types.Number{Value: l.Value + r.Value}, nil
	}
}

// evalBinaryDivide interprets a binary division e.g. `10 / 2`.
func evalBinaryDivide(left, right types.Type) (types.Number, error) {
	var zero types.Number
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return zero, err
	}

	result := types.Number{Value: l.Value / r.Value}
	return result, nil
}

// evalBinaryMultiply interprets a binary multiplication e.g. `5 * 3`.
func evalBinaryMultiply(left, right types.Type) (types.Number, error) {
	var zero types.Number
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return zero, err
	}

	result := types.Number{Value: l.Value * r.Value}
	return result, nil
}

// evalGreaterThan interprets e.g. `5 > 3`.
func evalGreaterThan(left, right types.Type) (*types.Bool, error) {
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return types.False, err
	}

	if l.Value > r.Value {
		return types.True, nil
	}
	return types.False, nil
}

// evalGreaterThanEq interprets e.g. `x >= y`.
func evalGreaterThanEq(left, right types.Type) (*types.Bool, error) {
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return types.False, err
	}

	if l.Value >= r.Value {
		return types.True, nil
	}
	return types.False, nil
}

// evalLessThan interprets `x < y`.
func evalLessThan(left, right types.Type) (*types.Bool, error) {
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return types.False, err
	}

	if l.Value < r.Value {
		return types.True, nil
	}
	return types.False, nil
}

// evalLessThanEq interprets e.g. `x <= y`.
func evalLessThanEq(left, right types.Type) (*types.Bool, error) {
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return types.False, err
	}

	if l.Value <= r.Value {
		return types.True, nil
	}
	return types.False, nil
}

// checkNumeric is a helper function to validate that the left and right operands of a binary
// expression are [types.Number], so that maths can be done on them.
//
// It returns the converted numbers and an error if the conversion could not occur.
func checkNumeric(left, right types.Type) (leftNumber, rightNumber types.Number, err error) {
	var zero types.Number
	leftNumber, ok := left.(types.Number)
	if !ok {
		return zero, zero, fmt.Errorf("left operand must be numeric, got %T", left)
	}

	rightNumber, ok = right.(types.Number)
	if !ok {
		return zero, zero, fmt.Errorf("right operand must be numeric, got %T", right)
	}

	return leftNumber, rightNumber, nil
}
