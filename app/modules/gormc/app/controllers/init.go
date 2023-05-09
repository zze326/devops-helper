package gormcontroller

import (
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		if revel.Config.BoolDefault("db.autoinit", true) {
			revel.InterceptMethod((*TxnController).Begin, revel.BEFORE)
			revel.InterceptMethod((*TxnController).Commit, revel.AFTER)
			revel.InterceptMethod((*TxnController).Rollback, revel.FINALLY)
			revel.InterceptMethod((*Controller).init, revel.BEFORE)
		}
	})
}
