package expr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Định nghĩa kiểu Node
type Node struct {
	Type      string  // Loại: "operator" (AND, OR, NOT, =, LIKE, +, -, ...), "operand" (Salary, MaxSalary, ...), "function" (LEN, ...)
	Value     string  // Giá trị: "AND", "LIKE", "+", "Salary", "LEN", ...
	Left      *Node   // Nút con trái (hoặc đối số của hàm)
	Right     *Node   // Nút con phải
	Arguments []*Node // Đối số của hàm
	Params    []string
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
		argsExpr := strings.TrimSpace(expr[idx+1 : len(expr)-1])
		argExpr, err := SplitFunctionExpression("f(" + argsExpr + ")")
		if err != nil {
			panic("error syntax: " + expr)
		}

		argsNode := []*Node{}
		for i := 0; i < len(argExpr); i++ {
			argNode := parseMathExpression(argExpr[i])
			argsNode = append(argsNode, argNode)

		}

		return &Node{
			Type:      "function",
			Value:     funcName,
			Left:      nil, // Đối số của hàm
			Right:     nil,
			Arguments: argsNode,
		}
	}

	// Nếu không có toán tử, đây là operand (biến hoặc giá trị)
	return &Node{Type: "operand", Value: expr}
}

var opsLogic = []string{"OR", "AND", "NOT", "||", "&&", "!"}

func ParseExpression(expr string) *Node {
	ret, params := parseExpression(expr, []string{})
	if len(params) > 0 {

		ret = rebuildNode(ret, params)
	}
	ret.Params = params
	return ret
}
func rebuildNode(node *Node, params []string) *Node {

	if node.Type == "operand" {
		if len(node.Value) > 2 && node.Value[0:2] == "@p" {
			strIndexOfParams := node.Value[2:]
			//check index can convert to int
			if indexOfParams, err := strconv.Atoi(strIndexOfParams); err == nil {
				if indexOfParams < len(params) {
					strP := params[indexOfParams]
					//escape ' in strP
					strP = strings.Replace(strP, "'", "''", -1)
					node.Value = strP

					node.Value = "'" + strP + "'"
				}
			}
		}

	}
	if node.Left != nil {
		node.Left = rebuildNode(node.Left, params)
	}
	if node.Right != nil {
		node.Right = rebuildNode(node.Right, params)
	}

	for i := 0; i < len(node.Arguments); i++ {
		if node.Arguments[i] != nil {
			node.Arguments[i] = rebuildNode(node.Arguments[i], params)
		}

	}
	return node

}

// Parse biểu thức logic
func parseExpression(expr string, params []string) (*Node, []string) {
	parseParams := ParseExprWithParam(expr)
	expr = parseParams.Expr
	expr = strings.TrimSpace(expr)

	// Xử lý dấu ngoặc: (expr)
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return parseExpression(expr[1:len(expr)-1], parseParams.Params)
	}
	for _, op := range opsLogic {
		if idx := findOperator(expr, op, len(op)); idx != -1 {
			leftExpr := strings.TrimSpace(expr[:idx])
			rightExpr := strings.TrimSpace(expr[idx+len(op):])
			if op == "!" {
				op = "NOT"
			}
			if op == "&&" {
				op = "AND"
			}
			if op == "||" {
				op = "OR"
			}
			return &Node{
				Type:  "operator",
				Value: op,
				Left:  ParseExpression(leftExpr),
				Right: ParseExpression(rightExpr),
			}, parseParams.Params
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
		{"==", 2},
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
			if comp.op == "==" {
				comp.op = "="
			}
			return &Node{
				Type:  "operator",
				Value: comp.op,
				Left:  parseMathExpression(leftExpr),
				Right: parseMathExpression(rightExpr),
			}, parseParams.Params
		}
	}

	// Nếu không có toán tử logic hoặc so sánh, parse như biểu thức toán học
	return parseMathExpression(expr), parseParams.Params
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
		//escape ' in node.Value

		return node.Value

	case "function":
		// Nếu là hàm, tái tạo cú pháp FunctionName(Argument)
		strArgs := []string{}
		for _, arg := range node.Arguments {
			strArgs = append(strArgs, ReconstructExpression(arg))
		}
		return fmt.Sprintf("%s(%s)", node.Value, strings.Join(strArgs, ","))

	case "operator":
		// Nếu là toán tử, tái tạo với các phần con
		switch node.Value {
		case "NOT", "!":
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
		case "&&":
			// Toán tử logic hoặc so sánh
			// Thêm dấu ngoặc nếu cần để bảo toàn độ ưu tiên
			node.Value = "AND"
			left := ReconstructExpression(node.Left)
			if node.Left.Type == "operator" && (node.Left.Value == "OR" || node.Left.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				left = fmt.Sprintf("(%s)", left)
			}
			right := ReconstructExpression(node.Right)
			if node.Right.Type == "operator" && (node.Right.Value == "OR" || node.Right.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				right = fmt.Sprintf("(%s)", right)
			}
			return fmt.Sprintf("%s %s %s", left, node.Value, right)
		case "||":
			// Toán tử logic hoặc so sánh
			// Thêm dấu ngoặc nếu cần để bảo toàn độ ưu tiên
			node.Value = "OR"
			left := ReconstructExpression(node.Left)
			if node.Left.Type == "operator" && (node.Left.Value == "OR" || node.Left.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				left = fmt.Sprintf("(%s)", left)
			}
			right := ReconstructExpression(node.Right)
			if node.Right.Type == "operator" && (node.Right.Value == "OR" || node.Right.Value == "AND") && node.Value != "AND" && node.Value != "OR" {
				right = fmt.Sprintf("(%s)", right)
			}
			return fmt.Sprintf("%s %s %s", left, node.Value, right)
		case "AND", "OR", "=", "LIKE", ">", "<", ">=", "<=", "<>", "!=", "==":
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
func SplitFunctionExpression(s string) ([]string, error) {
	if s == "" {
		return nil, fmt.Errorf("input string cannot be empty")
	}

	// Chuyển chuỗi thành rune để xử lý Unicode
	runes := []rune(s)
	if len(runes) < 2 {
		return nil, fmt.Errorf("input string is too short")
	}

	// Tìm vị trí của ngoặc mở đầu tiên
	openParenIndex := -1
	for i, r := range runes {
		if r == '(' {
			openParenIndex = i
			break
		}
	}
	if openParenIndex == -1 {
		return nil, fmt.Errorf("no opening parenthesis found")
	}

	// Tìm vị trí của ngoặc đóng cuối cùng
	parenDepth := 0
	closeParenIndex := -1
	for i, r := range runes {
		if r == '(' {
			parenDepth++
		} else if r == ')' {
			parenDepth--
			if parenDepth == 0 {
				closeParenIndex = i
				break
			}
		}
	}
	if closeParenIndex == -1 || parenDepth != 0 {
		return nil, fmt.Errorf("unmatched parentheses")
	}

	// Kiểm tra nếu không có tham số (ví dụ: "func()")
	if closeParenIndex-openParenIndex <= 1 {
		return []string{}, nil
	}

	// Tách các tham số trong phạm vi giữa openParenIndex và closeParenIndex
	params := []string{}
	start := openParenIndex + 1
	parenDepth = 0

	for i := openParenIndex + 1; i < closeParenIndex; i++ {
		r := runes[i]

		if r == '(' {
			parenDepth++
		} else if r == ')' {
			parenDepth--
			if parenDepth < 0 {
				return nil, fmt.Errorf("unmatched closing parenthesis at position %d", i)
			}
		} else if r == ',' && parenDepth == 0 {
			// Dấu phẩy ở mức ngoài cùng, tách tham số
			param := strings.TrimSpace(string(runes[start:i]))
			if param == "" {
				return nil, fmt.Errorf("empty parameter at position %d", start)
			}
			params = append(params, param)
			start = i + 1
		}
	}

	// Thêm tham số cuối cùng
	lastParam := strings.TrimSpace(string(runes[start:closeParenIndex]))
	if lastParam == "" {
		return nil, fmt.Errorf("empty parameter at position %d", start)
	}
	params = append(params, lastParam)

	return params, nil
}

type ExprWithParam struct {
	Expr   string
	Params []string
}

func ParseExprWithParam(e string) ExprWithParam {
	var result strings.Builder
	var params []string
	var i int
	paramIndex := 0

	for i < len(e) {
		// Tìm field name
		start := i
		for i < len(e) && ((e[i] >= 'a' && e[i] <= 'z') || (e[i] >= 'A' && e[i] <= 'Z') || (e[i] >= '0' && e[i] <= '9') || e[i] == '_') {
			i++
		}
		field := strings.TrimSpace(e[start:i])

		// Bỏ khoảng trắng
		for i < len(e) && e[i] == ' ' {
			i++
		}

		// Kiểm tra ==
		if i+1 < len(e) && e[i] == '=' && e[i+1] == '=' {
			i += 2
		} else {
			// Không khớp, copy ký tự thường vào
			result.WriteByte(e[start])
			i = start + 1
			continue
		}

		// Bỏ khoảng trắng
		for i < len(e) && e[i] == ' ' {
			i++
		}

		// Bắt đầu chuỗi value
		if i < len(e) && e[i] == '\'' {
			i++ // Bỏ dấu '

			valBuilder := strings.Builder{}

			for i < len(e) {
				if e[i] == '\\' && i+1 < len(e) {
					// Escape character
					valBuilder.WriteByte(e[i+1])
					i += 2
				} else if e[i] == '\'' {
					// Kết thúc chuỗi
					break
				} else {
					valBuilder.WriteByte(e[i])
					i++
				}
			}

			value := valBuilder.String()
			params = append(params, value)
			result.WriteString(fmt.Sprintf("%s=@p%d", field, paramIndex))
			paramIndex++
			i++ // Bỏ dấu kết thúc '
		} else {
			// Không phải chuỗi hợp lệ
			result.WriteByte(e[start])
			i = start + 1
		}

		// Thêm các ký tự còn lại (AND, OR, ...)
		for i < len(e) && (e[i] == ' ' || (e[i] >= 'A' && e[i] <= 'Z')) {
			result.WriteByte(e[i])
			i++
		}
	}

	return ExprWithParam{
		Expr:   result.String(),
		Params: params,
	}
}
func (node *Node) String() string {
	// convert to json with indentation
	jsonStr, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonStr)

}

type ISqlParserBase interface {
	ToAst(rawSql string) (*Node, error)
	ToSnakeCase(str string) string
	ReconstructExpression(node *Node) string
}
type ISqlParser interface {
	ISqlParserBase
	Conditional(rawCondition string) (string, error)
	ResolveAst(ast *Node) (*Node, error)
}

type SqlParserBase struct {
}

func (s SqlParserBase) ToSnakeCase(str string) string {
	return toSnakeCase(str)
}
func (s SqlParserBase) ToAst(rawSql string) (*Node, error) {
	return ParseExpression(rawSql), nil
}
func (s SqlParserBase) ReconstructExpression(node *Node) string {
	return ReconstructExpression(node)
}

//=======================================

func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	// Kiểm tra xem chuỗi có phải toàn chữ hoa (hoặc chữ hoa + số) không
	isAllUpper := true
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsNumber(r) && unicode.IsLetter(r) {
			isAllUpper = false
			return s
		}
	}

	// Nếu toàn chữ hoa, chỉ cần chuyển thành chữ thường
	if isAllUpper {
		return strings.ToLower(s)
	}

	var result strings.Builder
	runes := []rune(s)

	// Vị trí bắt đầu của chuỗi chữ hoa
	upperRunStart := -1

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsUpper(r) {
			if i == 0 {
				// Ký tự đầu tiên là chữ hoa, không thêm _
				result.WriteRune(unicode.ToLower(r))
				upperRunStart = i
			} else {
				// Kiểm tra ranh giới từ
				prevIsLower := unicode.IsLower(runes[i-1])
				nextIsLower := (i+1 < len(runes)) && unicode.IsLower(runes[i+1])

				if prevIsLower || (nextIsLower && upperRunStart != i-1) {
					// Thêm _ nếu trước đó là chữ thường hoặc đây là chữ hoa bắt đầu từ mới
					if result.Len() > 0 && result.String()[result.Len()-1] != '_' {
						result.WriteRune('_')
					}
				}
				result.WriteRune(unicode.ToLower(r))
				upperRunStart = i
			}
		} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			// Thay ký tự đặc biệt bằng dấu gạch dưới
			if result.Len() > 0 && result.String()[result.Len()-1] != '_' {
				result.WriteRune('_')
			}
			upperRunStart = -1
		} else {
			// Chữ thường hoặc số
			result.WriteRune(r)
			upperRunStart = -1
		}
	}

	// Loại bỏ dấu gạch dưới ở đầu và cuối, thay thế nhiều dấu gạch dưới liên tiếp bằng một
	snake := strings.Trim(result.String(), "_")
	snake = strings.ReplaceAll(snake, "__", "_")
	return snake
}
