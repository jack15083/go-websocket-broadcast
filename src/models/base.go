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

	return db, nil
}

func (baseModel) CloseDB(db *gorm.DB) {
	db.Close()
}
