package factory

import (
	"vngom/gormex/expr"
	_ "vngom/gormex/expr"
	"vngom/gormex/expr/exprpostgres"
	_ "vngom/gormex/expr/exprpostgres"
)

func NewExpr(driver string) expr.IExpr {
	switch driver {
	case "postgres":
		return exprpostgres.New()
	default:
		panic("Unsupported driver: " + driver)
	}
}
