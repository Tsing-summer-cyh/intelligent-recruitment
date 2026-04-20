package initialize

import (
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize/internal"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// GormSqlite 初始化 sqlite 数据库
func GormSqlite() *gorm.DB {
	s := global.GVA_CONFIG.Sqlite
	return initSqliteDatabase(s)
}

// GormSqliteByConfig 使用指定配置初始化 sqlite 数据库
func GormSqliteByConfig(s config.Sqlite) *gorm.DB {
	return initSqliteDatabase(s)
}

func initSqliteDatabase(s config.Sqlite) *gorm.DB {
	if s.Dbname == "" {
		return nil
	}
	if s.Path != "" {
		if err := os.MkdirAll(s.Path, os.ModePerm); err != nil {
			panic(err)
		}
	}

	general := s.GeneralDB
	if db, err := gorm.Open(sqlite.Open(s.Dsn()), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}
