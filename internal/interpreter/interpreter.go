// Package interpreter implements the tree walking interpreter.
package interpreter

import (
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/glox/internal/syntax/types"
)

// TODO(@FollowTheProcess): Currently it stops on the first error, we should report the error
// to the user with position info and then move on to the next expression. Maybe need an interpreter
// struct with fields tracking errors etc.

// TODO(@FollowTheProcess): Make the errors a lot better with position info and highlighting etc.

// Interpreter is the glox interpreter.
type Interpreter struct {
	env *Environment // Global interpreter environment
}

// New returns a new interpreter.
func New() Interpreter {
	return Interpreter{
		env: NewEnvironment("globals", nil),
	}
}

// Eval evaluates an AST node.
func (i Interpreter) Eval(node ast.Node) (types.Type, error) {
	switch node := node.(type) {
	case ast.Program:
		return i.evalStatements(node.Statements)
	case ast.ExpressionStatement:
		return i.Eval(node.Value)
	case ast.DeclarationStatement:
		return i.Eval(node.Value)
	case ast.GroupedExpression:
		return i.Eval(node.Value)
	case ast.UnaryExpression:
		return i.evalUnaryExpression(node)
	case ast.BinaryExpression:
		return i.evalBinaryExpression(node)
	case ast.VarDeclaration:
		return nil, i.evalVarDeclaration(node)
	case ast.Number:
		return i.evalNumber(node), nil
	case ast.Bool:
		return i.evalBool(node), nil
	case ast.String:
		return i.evalString(node), nil
	case ast.Ident:
		return i.evalIdent(node)
	default:
		return nil, fmt.Errorf("unhandled ast.Node in Eval: %T", node)
	}
}

// evalStatements iterates through all the statements in the program, evaluating each
// and returning the final type.
func (i Interpreter) evalStatements(statements []ast.Statement) (types.Type, error) {
	var result types.Type
	var err error
	for _, statement := range statements {
		result, err = i.Eval(statement)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// evalBool interprets a bool AST node, returning either the True or False singleton.
func (i Interpreter) evalBool(node ast.Bool) *types.Bool {
	if node.Value {
		return types.True
	}
	return types.False
}

// evalNumber interprets a numeric AST node, returning a types.Number.
func (i Interpreter) evalNumber(node ast.Number) types.Number {
	return types.Number{Value: node.Value}
}

// evalString interprets a string AST node, returning a types.String.
func (i Interpreter) evalString(node ast.String) types.String {
	return types.String{Value: node.Value}
}

// evalUnaryExpression interprets a unary expression like `-5` or `!true`.
func (i Interpreter) evalUnaryExpression(node ast.UnaryExpression) (types.Type, error) {
	operand, err := i.Eval(node.Value)
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
func (i Interpreter) evalBinaryExpression(node ast.BinaryExpression) (types.Type, error) {
	left, err := i.Eval(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Eval(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Op.Kind {
	case token.Minus:
		return i.evalBinarySubtract(left, right)
	case token.Plus:
		return i.evalBinaryAdd(left, right)
	case token.ForwardSlash:
		return i.evalBinaryDivide(left, right)
	case token.Star:
		return i.evalBinaryMultiply(left, right)
	case token.GreaterThan:
		return i.evalGreaterThan(left, right)
	case token.GreaterThanEq:
		return i.evalGreaterThanEq(left, right)
	case token.LessThan:
		return i.evalLessThan(left, right)
	case token.LessThanEq:
		return i.evalLessThanEq(left, right)
	case token.DoubleEq:
		return i.evalEqual(left, right), nil
	case token.BangEq:
		return i.evalNotEqual(left, right), nil
	default:
		return nil, fmt.Errorf("unsupported binary operator: %s", node.Op.Kind.Lexeme())
	}
}

// evalBinarySubtract interprets a binary subtraction e.g. `5 - 3`.
func (i Interpreter) evalBinarySubtract(left, right types.Type) (types.Number, error) {
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
func (i Interpreter) evalBinaryAdd(left, right types.Type) (types.Type, error) {
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
func (i Interpreter) evalBinaryDivide(left, right types.Type) (types.Number, error) {
	var zero types.Number
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return zero, err
	}

	result := types.Number{Value: l.Value / r.Value}
	return result, nil
}

// evalBinaryMultiply interprets a binary multiplication e.g. `5 * 3`.
func (i Interpreter) evalBinaryMultiply(left, right types.Type) (types.Number, error) {
	var zero types.Number
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return zero, err
	}

	result := types.Number{Value: l.Value * r.Value}
	return result, nil
}

// evalGreaterThan interprets e.g. `5 > 3`.
func (i Interpreter) evalGreaterThan(left, right types.Type) (*types.Bool, error) {
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
func (i Interpreter) evalGreaterThanEq(left, right types.Type) (*types.Bool, error) {
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
func (i Interpreter) evalLessThan(left, right types.Type) (*types.Bool, error) {
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
func (i Interpreter) evalLessThanEq(left, right types.Type) (*types.Bool, error) {
	l, r, err := checkNumeric(left, right)
	if err != nil {
		return types.False, err
	}

	if l.Value <= r.Value {
		return types.True, nil
	}
	return types.False, nil
}

// evalEqual interprets `x == y`.
func (i Interpreter) evalEqual(left, right types.Type) *types.Bool {
	if types.Equal(left, right) {
		return types.True
	}
	return types.False
}

// evalNotEqual interprets `x != y`.
func (i Interpreter) evalNotEqual(left, right types.Type) *types.Bool {
	if types.Equal(left, right) {
		return types.False
	}
	return types.True
}

// evalVarDeclaration evaluates `var <ident> = <expr>`.
func (i Interpreter) evalVarDeclaration(node ast.VarDeclaration) error {
	if node.Value == nil {
		// We're just defining the variable, no expression
		return i.env.Define(node.Ident.Name, nil)
	}

	value, err := i.Eval(node.Value)
	if err != nil {
		return err
	}

	return i.env.Define(node.Ident.Name, value)
}

// evalIdent evaluates an ident expression.
func (i Interpreter) evalIdent(node ast.Ident) (types.Type, error) {
	v, defined := i.env.Get(node.Name)
	if !defined {
		return nil, fmt.Errorf("use of undefined variable %q", node.Name)
	}

	return v, nil
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
