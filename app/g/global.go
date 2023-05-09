package g

import (
	"github.com/casbin/casbin/v2"
	"github.com/revel/revel/logger"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Logger logger.MultiLogger
var Enforcer *casbin.Enforcer
