package models

import (
	"time"

	"../config"
)

type PushMessageModel struct {
	ID         int64
	Title      string
	Content    string
	Options    string
	MsgType    int
	UserIds    string
	SenderId   int64
	SenderName string
	CreateTime string
}

func (PushMessageModel) TableName() string {
	return "xhx_push_message"
}

func (PushMessageModel) Create(m PushMessageModel) int64 {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return 0
	}
	defer db.Close()

	m.CreateTime = time.Now().Format(config.TIMESTAMP_FORMAT)

	db.Create(&m)

	return m.ID
}
