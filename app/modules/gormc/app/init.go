package gormdb

// # Database config
// db.driver=sqlite3 # mysql, postgres, sqlite3
// db.host=localhost  # Use dbhost  /tmp/app.db is your driver is sqlite
// db.port=dbport
// db.user=dbuser
// db.name=dbname
// db.password=dbpassword
// db.singulartable=false # default=false

import (
	"fmt"
	"github.com/zze326/devops-helper/app/g"
	"gorm.io/gorm/logger"
	"strings"

	"gorm.io/gorm"

	"github.com/revel/revel"
	"gorm.io/driver/mysql"
)

// DB Gorm.
var (
	DB      *gorm.DB
	gormLog = revel.AppLog
)

func init() {
	revel.RegisterModuleInit(func(m *revel.Module) {
		gormLog = m.Log
		g.Logger = m.Log
	})
}

// InitDB database.
func OpenDB(params DbInfo) {
	dbInfo := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", params.DbUser, params.DbPassword, params.DbHost, params.DbPort, params.DbName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbInfo, // DSN data source name
		DefaultStringSize:         256,    // string 类型字段的默认长度
		DisableDatetimePrecision:  true,   // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,   // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,   // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,  // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		// 禁止创建物理外键
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		gormLog.Fatal("GORM: open db failed", "error", err)
	}
	DB = db
	sqlDB, err := db.DB()
	if err != nil {
		gormLog.Fatal("GORM: get db failed", "error", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	gormLog.Infof("GORM: connected to %s", strings.ReplaceAll(dbInfo, params.DbPassword, "********"))
}

type DbInfo struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
}

func InitDB() {
	params := DbInfo{}
	params.DbHost = revel.Config.StringDefault("db.host", "localhost")
	params.DbPort = revel.Config.IntDefault("db.port", 3306)
	params.DbUser = revel.Config.StringDefault("db.user", "default")
	params.DbPassword = revel.Config.StringDefault("db.password", "")
	params.DbName = revel.Config.StringDefault("db.name", "default")

	OpenDB(params)
}

func AutoMigrate(dstModels ...interface{}) {
	err := DB.AutoMigrate(dstModels...)
	if err != nil {
		gormLog.Fatal("GORM: migrate db failed", "error", err)
	}
}
