package sql_ast_parse_test

import (
	"fmt"
	"testing"
	"vngom/repo/sql_ast_parse"
)

func TestParse(t *testing.T) {
	// Biểu thức phức tạp
	expression := "(Salary + Bonus) > MaxSalary AND LEN(FirstName) > 5"
	expression = "(Salary + Bonus) > MaxSalary and LEN(FirstName) > ? "
	// Parse thành AST
	ast := sql_ast_parse.ParseExpression(expression)
	rc := sql_ast_parse.ReconstructExpression(ast)
	fmt.Println("Reconstructed expression:", rc)
	if ast != nil {
		fmt.Println("AST Structure:")
		sql_ast_parse.TraverseAST(ast, 0)
	} else {
		fmt.Println("Invalid expression")
	}
}
