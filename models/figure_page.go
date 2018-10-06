package models

import (
	"github.com/jinzhu/gorm"
)

// FigurePage for page template with query of a display page
type FigurePage struct {
	ID   uint   `gorm:"primary_key;auto_increment"`
	Data string `gorm:"column:data;type:jsonb"`
}

// GetFigurePage get figure page from DB
func GetFigurePage(db *gorm.DB, id string) *FigurePage {
	ret := new(FigurePage)
	err := db.Where("data->>'id' = ?", id).First(ret).Error
	if err != nil {
		return nil
	}
	return ret
}
