package models

import (
	"../core"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var _db map[string]*gorm.DB

func Init() {
	_db = make(map[string]*gorm.DB)
	dBConf := core.Config.Database
	for key, conf := range dBConf {
		db, err := gorm.Open(conf.DriverName, conf.DataSourceName)
		if err != nil {
			panic("连接数据库失败, error=" + err.Error())
		}

		//设置最大空闲连接数
		db.DB().SetMaxIdleConns(conf.MaxIdleNum)
		//设置最大连接数
		db.DB().SetMaxOpenConns(conf.MaxOpenNum)
		// 开启 Logger, 以展示详细的db日志
		//db.LogMode(core.Config.Logger.Debug)
		dbLog := core.Config.Logger.New("db_" + key)
		db.SetLogger(dbLog)
		_db[key] = db
	}
}

func GetDB(DBName string) *gorm.DB {
	return _db[DBName]
}

func CloseDB(DBName string) {
	_db[DBName].Close()
}
