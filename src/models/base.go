package models

import (
	"../core"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var BaseModel baseModel

type baseModel struct{}

func (baseModel) ConnectDB(DBName string) (db *gorm.DB, err error) {
	dBConf := core.Config.Database[DBName]
	db, connectErr := gorm.Open(dBConf.DriverName, dBConf.DataSourceName)
	if connectErr != nil {
		log.Error("error connect db." + connectErr.Error())
		return nil, connectErr
	}

	//设置最大空闲连接数
	db.DB().SetMaxIdleConns(dBConf.MaxIdleNum)
	//设置最大连接数
	db.DB().SetMaxOpenConns(dBConf.MaxOpenNum)
	// 开启 Logger, 以展示详细的db日志
	//db.LogMode(core.Config.Logger.Debug)
	dbLog := core.Config.Logger.New("db_error")
	db.SetLogger(dbLog)

	return db, nil
}

func (baseModel) CloseDB(db *gorm.DB) {
	db.Close()
}
