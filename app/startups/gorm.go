package startups

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	"github.com/zze326/devops-helper/app/models/orm"
	gormdb "github.com/zze326/devops-helper/app/modules/gormc/app"
)

func initDB() {
	if revel.Config.BoolDefault("db.autoinit", true) {
		gormdb.InitDB()
	}

	if revel.Config.BoolDefault("db.automigrate", true) {
		gormdb.AutoMigrate(orm.ModelTypesToMigrate()...)
	}

	g.DB = gormdb.DB
}
