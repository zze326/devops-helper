package gormcontroller

import (
	"github.com/revel/revel"
	gormdb "github.com/zze326/devops-helper/app/modules/gormc/app"
	"gorm.io/gorm"
)

// Controller is a Revel controller with a pointer to the opened database.
type Controller struct {
	*revel.Controller
	DB *gorm.DB
}

func (c *Controller) init() revel.Result {
	c.DB = gormdb.DB
	return nil
}
