package models

import (
	"github.com/jinzhu/gorm"
)

// Figure for a figure template with parse query
type Figure struct {
	ID   uint   `gorm:"primary_key;auto_increment"`
	Data string `gorm:"column:data;type:jsonb"`
}

// GetFigure get figure from DB
func GetFigure(db *gorm.DB, id string) *Figure {
	ret := new(Figure)
	err := db.Where("data->>'id' = ?", id).First(ret).Error
	if err != nil {
		return nil
	}
	return ret
}
