package exprpostgres

import (
	"strings"
	"sync"
	"vngom/gormex/expr"
)

type ExprPostgres struct {
	expr.ISqlParserBase
}

func (ep *ExprPostgres) Conditional(rawCondition string) (string, error) {
	ast := expr.ParseExpression(rawCondition)
	newAst, err := ep.ResolveAst(ast)
	if err != nil {
		return "", err
	}

	return ep.ReconstructExpression(newAst), nil
}
func (ep *ExprPostgres) ResolveAst(ast *expr.Node) (*expr.Node, error) {
	return resolveAst(ast, ep)
}

var (
	insOfExprPostgres *ExprPostgres
	once              sync.Once
)

func New() expr.ISqlParser {
	once.Do(func() {
		insOfExprPostgres = &ExprPostgres{
			expr.SqlParserBase{},
		}
	})
	return insOfExprPostgres
}
func resolveAst(ast *expr.Node, ep expr.ISqlParserBase) (*expr.Node, error) {
	if ast == nil {
		return nil, nil
	}

	retAst := ast
	switch ast.Type {
	case "operand":
		retAst.Value = ep.ToSnakeCase(ast.Value)
		l, le := resolveAst(ast.Left, ep)
		if le != nil {
			return nil, le
		}
		retAst.Left = l
		r, re := resolveAst(ast.Right, ep)
		if re != nil {
			return nil, re
		}
		retAst.Right = r
		return retAst, nil
	case "function":
		return resolvePgFunction(ast, ep)

	default:
		for i, arg := range ast.Arguments {
			argAst, err := resolveAst(arg, ep)
			if err != nil {
				return nil, err
			}
			ast.Arguments[i] = argAst
		}
		rLeft, err := resolveAst(ast.Left, ep)
		if err != nil {
			return nil, err
		}

		retRight, err := resolveAst(ast.Right, ep)
		if err != nil {
			return nil, err
		}
		retAst.Left = rLeft
		retAst.Right = retRight
		return retAst, nil
	}

}

// this function will transform the function to postgres function
func resolvePgFunction(ast *expr.Node, ep expr.ISqlParserBase) (*expr.Node, error) {
	retNode := ast

	funcName := strings.ToLower(ast.Value)
	switch funcName {
	case "len":

		retNode.Value = "LENGTH"
		for i, arg := range ast.Arguments {
			argAst, err := resolveAst(arg, ep)
			if err != nil {
				return nil, err
			}
			retNode.Arguments[i] = argAst
		}
		return retNode, nil
	case "year", "month", "day", "hour", "minute", "second":
		retNode.Value = "date_part"
		oldArg := retNode.Arguments[0]
		oldArg.Value = ep.ToSnakeCase(oldArg.Value)
		retNode.Arguments = []*expr.Node{
			{
				Type:  "operand",
				Value: "'" + funcName + "'",
			},
			oldArg,
		}

		return retNode, nil
	default:
		for _, arg := range ast.Arguments {
			argAst, err := resolveAst(arg, ep)
			if err != nil {
				return nil, err
			}
			retNode.Arguments = append(retNode.Arguments, argAst)
		}
		return retNode, nil
	}
}
