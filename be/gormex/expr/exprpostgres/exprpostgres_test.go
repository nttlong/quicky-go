package exprpostgres_test

import (
	"strings"
	"testing"

	"vngom/gormex/expr/exprpostgres"

	"github.com/stretchr/testify/assert"
)

var testData = []string{
	"year(ID)->date_part('year',id)",
	"month(ID)->date_part('month',id)",
	"day(ID)->date_part('day',id)",
	"hour(ID)->date_part('hour',id)",
	"minute(ID)->date_part('minute',id)",
	"second(ID)->date_part('second',id)",
	"ID->id",
}

func TestParseConditional(t *testing.T) {
	parser := exprpostgres.New()
	for _, test := range testData {
		input := strings.Split(test, "->")[0]
		ouput := strings.Split(test, "->")[1]
		expr, err := parser.Conditional(input)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, ouput, expr)
	}

}
