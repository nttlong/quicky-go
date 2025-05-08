package compiler_test

import (
	"fmt"
	"strings"
	"testing"
	"vngom/gormex/expr/compiler"

	"github.com/stretchr/testify/assert"
)

var testList []string = []string{
	"(year(CreatedOn) == ?) && (month(CreatedOn) == ?)->(year(CreatedOn) == ?) && (month(CreatedOn) == ?)",
	"(Code==? and Price<=?) or len(name)==?->(Code == ? and Price <= ?) or len(name) == ?",
	"(Code==? and Price<=?) || len(name)==?->(Code == ? and Price <= ?) || len(name) == ?",
}

func TestTree(t *testing.T) {
	resolver := func(n *compiler.SimpleExprTree) error {
		// if n.Nt == "func" {
		// 	n.V = "pg." + n.V

		// }
		return nil
	}
	for _, s := range testList {
		Input := strings.Split(s, "->")[0]
		Output := strings.Split(s, "->")[1]
		fx, err := compiler.ParseExpr(Input)
		fmt.Print(fx.String())
		if err != nil {
			t.Error(err)
		}
		r, err := compiler.Resolve(fx, resolver)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, Output, r)
	}
}
