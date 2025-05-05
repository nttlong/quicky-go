package sql_ast_parse

import (
	"fmt"
	"strings"
)

// Định nghĩa kiểu Node
type Node struct {
	Type  string // Loại: "operator" (AND, OR, NOT, =, LIKE, +, -, ...), "operand" (Salary, MaxSalary, ...), "function" (LEN, ...)
	Value string // Giá trị: "AND", "LIKE", "+", "Salary", "LEN", ...
	Left  *Node  // Nút con trái (hoặc đối số của hàm)
	Right *Node  // Nút con phải
}

// Tìm vị trí toán tử, bỏ qua phần trong dấu ngoặc
func findOperator(Expr string, op string, opLen int) int {
	expr := strings.ToUpper(Expr)
	parenCount := 0
	for i := 0; i < len(expr); i++ {
		if expr[i] == '(' {
			parenCount++
		} else if expr[i] == ')' {
			parenCount--
		} else if parenCount == 0 && i+opLen <= len(expr) && strings.HasPrefix(expr[i:], op) {
			return i
		}
	}
	return -1
}

// Parse biểu thức toán học (+, -, *, /)
func parseMathExpression(expr string) *Node {
	expr = strings.TrimSpace(expr)

	// Xử lý dấu ngoặc: (expr)
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return parseMathExpression(expr[1 : len(expr)-1])
	}

	// Tìm toán tử + hoặc - (ưu tiên thấp hơn *, /)
	for _, op := range []struct {
		op  string
		len int
	}{
		{"+", 1},
		{"-", 1},
	} {
		if idx := findOperator(expr, op.op, op.len); idx != -1 {
			leftExpr := strings.TrimSpace(expr[:idx])
			rightExpr := strings.TrimSpace(expr[idx+op.len:])
			return &Node{
				Type:  "operator",
				Value: op.op,
				Left:  parseMathExpression(leftExpr),
				Right: parseMathExpression(rightExpr),
			}
		}
	}

	// Tìm toán tử * hoặc / (ưu tiên cao hơn +, -)
	for _, op := range []struct {
		op  string
		len int
	}{
		{"*", 1},
		{"/", 1},
	} {
		if idx := findOperator(expr, op.op, op.len); idx != -1 {
			leftExpr := strings.TrimSpace(expr[:idx])
			rightExpr := strings.TrimSpace(expr[idx+op.len:])
			return &Node{
				Type:  "operator",
				Value: op.op,
				Left:  parseMathExpression(leftExpr),
				Right: parseMathExpression(rightExpr),
			}
		}
	}

	// Nếu không có toán tử, kiểm tra xem có phải hàm không (ví dụ: LEN(FirstName))
	if idx := strings.Index(expr, "("); idx != -1 && strings.HasSuffix(expr, ")") {
		funcName := strings.TrimSpace(expr[:idx])
		argExpr := strings.TrimSpace(expr[idx+1 : len(expr)-1])
		return &Node{
			Type:  "function",
			Value: funcName,
			Left:  parseMathExpression(argExpr), // Đối số của hàm
			Right: nil,
		}
	}

	// Nếu không có toán tử, đây là operand (biến hoặc giá trị)
	return &Node{Type: "operand", Value: expr}
}

// Parse biểu thức logic
func ParseExpression(expr string) *Node {
	expr = strings.TrimSpace(expr)

	// Xử lý dấu ngoặc: (expr)
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return ParseExpression(expr[1 : len(expr)-1])
	}

	// Tìm toán tử OR (ưu tiên thấp nhất)
	if idx := findOperator(expr, "OR", 2); idx != -1 {
		leftExpr := strings.TrimSpace(expr[:idx])
		rightExpr := strings.TrimSpace(expr[idx+2:])
		return &Node{
			Type:  "operator",
			Value: "OR",
			Left:  ParseExpression(leftExpr),
			Right: ParseExpression(rightExpr),
		}
	}

	// Tìm toán tử AND (ưu tiên cao hơn OR)
	if idx := findOperator(expr, "AND", 3); idx != -1 {
		leftExpr := strings.TrimSpace(expr[:idx])
		rightExpr := strings.TrimSpace(expr[idx+3:])
		return &Node{
			Type:  "operator",
			Value: "AND",
			Left:  ParseExpression(leftExpr),
			Right: ParseExpression(rightExpr),
		}
	}

	// Tìm toán tử NOT (ưu tiên cao hơn AND)
	if strings.HasPrefix(expr, "NOT ") {
		condition := strings.TrimSpace(expr[4:])
		return &Node{
			Type:  "operator",
			Value: "NOT",
			Left:  ParseExpression(condition),
			Right: nil, // NOT là toán tử một ngôi
		}
	}

	// Tìm các phép so sánh: =, LIKE, >, <, >=, <=, <>, !=
	comparators := []struct {
		op    string
		count int
	}{
		{"LIKE", 4},
		{"<>", 2},
		{"!=", 2},
		{">=", 2},
		{"<=", 2},
		{"=", 1},
		{">", 1},
		{"<", 1},
	}

	for _, comp := range comparators {
		if idx := findOperator(expr, comp.op, comp.count); idx != -1 {
			// Phân tách thành hai phần: trước và sau toán tử so sánh
			leftExpr := strings.TrimSpace(expr[:idx])
			rightExpr := strings.TrimSpace(expr[idx+comp.count:])

			// Parse từng phần bên trái và phải (có thể chứa phép toán học)
			return &Node{
				Type:  "operator",
				Value: comp.op,
				Left:  parseMathExpression(leftExpr),
				Right: parseMathExpression(rightExpr),
			}
		}
	}

	// Nếu không có toán tử logic hoặc so sánh, parse như biểu thức toán học
	return parseMathExpression(expr)
}

// Duyệt và in toàn bộ AST
func TraverseAST(node *Node, level int) {
	if node == nil {
		return
	}
	fmt.Printf("%sType: %s, Value: %s\n", strings.Repeat("  ", level), node.Type, node.Value)
	TraverseAST(node.Left, level+1)
	TraverseAST(node.Right, level+1)
}
func ReconstructExpression(node *Node) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case "operand":
		// Nếu là operand, trả về giá trị gốc
		return node.Value

	case "function":
		// Nếu là hàm, tái tạo cú pháp FunctionName(Argument)
		return fmt.Sprintf("%s(%s)", node.Value, ReconstructExpression(node.Left))

	case "operator":
		// Nếu là toán tử, tái tạo với các phần con
		switch node.Value {
		case "NOT":
			// Toán tử một ngôi NOT
			return fmt.Sprintf("NOT %s", ReconstructExpression(node.Left))

		case "+", "-", "*", "/":
			// Toán tử toán học (+, -, *, /)
			// Thêm dấu ngoặc để bảo toàn độ ưu tiên
			left := ReconstructExpression(node.Left)
			if node.Left.Type == "operator" && (node.Left.Value == "+" || node.Left.Value == "-") && (node.Value == "*" || node.Value == "/") {
				left = fmt.Sprintf("(%s)", left)
			}
			right := ReconstructExpression(node.Right)
			if node.Right.Type == "operator" && (node.Right.Value == "+" || node.Right.Value == "-") && (node.Value == "*" || node.Value == "/") {
				right = fmt.Sprintf("(%s)", right)
			}
			return fmt.Sprintf("%s %s %s", left, node.Value, right)

		case "AND", "OR", "=", "LIKE", ">", "<", ">=", "<=", "<>", "!=":
			// Toán tử logic hoặc so sánh
			// Thêm dấu ngoặc nếu cần để bảo toàn độ ưu tiên
			left := ReconstructExpression(node.Left)
			if node.Left.Type == "operator" && (node.Left.Value == "OR" || node.Left.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				left = fmt.Sprintf("(%s)", left)
			}
			right := ReconstructExpression(node.Right)
			if node.Right.Type == "operator" && (node.Right.Value == "OR" || node.Right.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				right = fmt.Sprintf("(%s)", right)
			}
			return fmt.Sprintf("%s %s %s", left, node.Value, right)

		default:
			return fmt.Sprintf("unknown_operator(%s)", node.Value)
		}

	default:
		return fmt.Sprintf("unknown_node(%s)", node.Value)
	}
}
