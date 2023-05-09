package gormcontroller

import (
	"database/sql"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	gormdb "github.com/zze326/devops-helper/app/modules/gormc/app"
	"gorm.io/gorm"
)

// TxnController is a Revel controller with database transaction support (begin, commit and rollback).
type TxnController struct {
	*revel.Controller
	Txn *gorm.DB
}

// Begin begins a DB transaction.
func (c *TxnController) Begin() revel.Result {
	txn := gormdb.DB.Begin()
	if txn.Error != nil {
		g.Logger.Panic("Transaction begine error", "error", txn.Error)
	}

	c.Txn = txn
	return nil
}

// Commit commits the database transition.
func (c *TxnController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}

	c.Txn.Commit()
	if c.Txn.Error != nil && c.Txn.Error != sql.ErrTxDone {
		g.Logger.Panic("Transaction commit error", "error", c.Txn.Error)
	}

	c.Txn = nil
	return nil
}

// Rollback rolls back the transaction (eg. after a panic).
func (c *TxnController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}

	c.Txn.Rollback()
	if c.Txn.Error != nil && c.Txn.Error != sql.ErrTxDone {
		g.Logger.Panic("Transaction rollback error", "error", c.Txn.Error)
	}

	c.Txn = nil
	return nil
}
