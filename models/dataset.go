package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Dataset represents a product data collection
type Dataset struct {
	gorm.Model
	IndexUpdatedAt           time.Time `gorm:"column:index_updated_at"`
	Name                     string    `gorm:"column:name;not null"`
	DisplayName              string    `gorm:"column:display_name"`
	Icon                     string    `gorm:"column:icon"`
	FigureSet                string    `gorm:"column:figure_set;type:jsonb"`
	CoreIndexName            string    `gorm:"column:core_index_name"`
	CoreIndexQuery           string    `gorm:"column:core_index_query;type:jsonb"`
	CoreIndexUpdateFrequency string    `gorm:"column:core_index_update_frequency"`
	IndexChange              float64   `gorm:"column:index_change;type:decimal(8,2)"`
}

// GetAllDatasets return all datasets in DB
func GetAllDatasets(db *gorm.DB) []Dataset {
	var ret []Dataset
	db.Find(&ret)
	return ret
}
