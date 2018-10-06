//  Created by paincompiler on 17/01/2018

package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"encoding/json"
)

type SearchHistoryDBOp struct {
	DB *gorm.DB
}

type SearchHistory struct {
	ID           uint      `gorm:"primary_key;auto_increment"`
	UserID       uint      `gorm:"column:user_id"`
	DatasetID    string    `gorm:"column:dataset_id"`
	FigurePageID string    `gorm:"column:figure_page_id"`
	TimeStamp    time.Time `gorm:"column:timestamp"`
}

func (SearchHistory) TableName() string {
	return "search_history"
}

func (r *SearchHistoryDBOp) GetLatestN(userID uint, n int) []SearchHistory {
	var ret []SearchHistory
	r.DB.Debug().Limit(n).Where("user_id = ?", userID).Order("id desc").Find(&ret)
	return ret
}

func (r *SearchHistoryDBOp) SaveSearchHistory(user_id uint, datasetID string, figurePageID string) {
	searchHistory := SearchHistory{
		UserID:       user_id,
		DatasetID:    datasetID,
		FigurePageID: figurePageID,
		TimeStamp:    time.Now(),
	}
	err := r.DB.Create(&searchHistory).Error
	if err != nil {
		s, _ := json.Marshal(searchHistory)
		fmt.Println("Failed to save search history[%s]!", string(s))
	}
}
