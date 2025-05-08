package expr_test

import (
	"fmt"
	"strings"
	"testing"

	"vngom/gormex/expr"

	"github.com/stretchr/testify/assert"
)

var exprlst = []string{
	"((year(join_date)=?) AND (month(join_date)=?)) && (a==b)->year(join_date) = ? AND month(join_date) OR (a=b) = ?",

	"(year(join_date)=?) AND (month(join_date)=?)->year(join_date) = ? AND month(join_date) = ?",
	"(year(join_date)=?) && (month(join_date)=?)->year(join_date) = ? AND month(join_date) = ?",

	"code=='abc\\' AND cc'->code = 'abc'' AND cc'",
	"code=='abc AND cc'->code = 'abc AND cc'",
	"code=='abc OR cc'->code = 'abc OR cc'",
	"code=='abc NOT cc'->code = 'abc NOT cc'",
	// "not (a&&b)->NOT a = b",
	"a&&b->a AND b",
	"a||b->a OR b",

	"a==b->a = b",
	"year(join_date)=2025->year(join_date) = 2025",
	"1+2*3->1 + 2 * 3",
	"1+2*3+4->1 + 2 * 3 + 4",
	"1+2*3+4-5->1 + 2 * 3 + 4 - 5",
	"1+2*3+4-5/6->1 + 2 * 3 + 4 - 5 / 6",
	"1+2*3+4-5/6*7->1 + 2 * 3 + 4 - 5 / 6 * 7",
	"year(join_date)=2025->year(join_date) = 2025",
	"sum(a,b,c)->sum(a,b,c)",
	"date_part('year',birth_day)=?->date_part('year',birth_day) = ?",
	"left(first_name+''+last_name,3)='abc'->left(first_name + '' + last_name,3) = 'abc'",
	"left(first_name+' '+last_name,3) like 'abc'->left(first_name + ' ' + last_name,3) LIKE 'abc'",
}

func TestSnakeCase(t *testing.T) {
	testList := []string{
		"UserName->user_name",
		"Email->email",
		"FirstName->first_name",
		"Code->code",
		"FullName->full_name",
		"Id->id",
		"UserId->user_id",
		"ID->id",
		"ABC->abc",
		"ABCD->abcd",
	}
	for _, x := range testList {
		input := strings.Split(x, "->")[0]
		ouput := strings.Split(x, "->")[1]
		fx := &expr.SqlParserBase{}
		r := fx.ToSnakeCase(input)
		assert.Equal(t, r, ouput)

	}

}
func TestParseExprWithParam(t *testing.T) {
	ex1 := expr.ParseExprWithParam("code=='abc\\' AND cc'")
	fmt.Println(ex1.Expr) // "code=@p0"
	assert.Equal(t, "abc' AND cc", ex1.Params[0])
	fmt.Println(ex1.Params) // ["abc AND cc"]

	ex2 := expr.ParseExprWithParam("code=='abc OR cc' AND name=='def'")
	fmt.Println(ex2.Expr)   // "code=@p0 AND name=@p1"
	fmt.Println(ex2.Params) // ["abc OR cc", "def"]
}
func TestExpr(t *testing.T) {
	for _, x := range exprlst {
		input := strings.Split(x, "->")[0]
		ouput := strings.Split(x, "->")[1]
		t.Log(x)
		e := expr.ParseExpression(input)
		r := expr.ReconstructExpression(e)
		//r = cleanSpaces(r)
		fmt.Println(e)
		fmt.Println(e.Value, "->", r, " ", x)
		assert.Equal(t, ouput, r)
	}
}
func TestExprParseWithFunction(t *testing.T) {
	x := "left(concat(first_name,' ',last_name),3)='abc'"
	e := expr.ParseExpression(x)
	r := expr.ReconstructExpression(e)
	// r = cleanSpaces(r)
	fmt.Println(e)
	t.Log(r)
	// fmt.Println(e.Value, "->", r, " ", x)
	// assert.Equal(t, x, r)
}
func TestSplitFunctionExperssion(t *testing.T) {
	st, e := expr.SplitFunctionExpression("f(concat(first_name,' ',last_name),1)")
	assert.Equal(t, e, nil)
	assert.Equal(t, st[0], "concat(first_name,' ',last_name)")
	assert.Equal(t, st[1], "1")
	st, e = expr.SplitFunctionExpression("f(first_name,' ',last_name)")
	assert.Equal(t, e, nil)
	assert.Equal(t, st[0], "first_name")
	assert.Equal(t, st[1], "' '")
	assert.Equal(t, st[2], "last_name")
}
