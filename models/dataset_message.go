package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DatasetMessage for messages of product data collection
type DatasetMessage struct {
	ID        uint      `gorm:"primary_key;auto_increment"`
	DatasetID uint      `gorm:"column:dataset_id;not null"`
	Msg       string    `gorm:"column:msg"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// GetDatasetMessages fetches messages from DB
func GetDatasetMessages(db *gorm.DB, datasetID uint) []DatasetMessage {
	msgs := make([]DatasetMessage, 0)
	db.Where(DatasetMessage{DatasetID: datasetID}).Find(&msgs)
	return msgs
}
