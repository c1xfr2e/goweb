package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	"encoding/json"
)

type RequestDBOp struct {
	DB *gorm.DB
}

type Request struct {
	ID         uint      `gorm:"primary_key;auto_increment"`
	UserID     uint      `gorm:"column:user_id"`
	Content    string    `gorm:"column:content"`
	SourceType string    `gorm:"column:source_type"`
	Extra      string    `gorm:"column:extra"`
	TimeStamp  time.Time `gorm:"column:timestamp"`
}

func (r *RequestDBOp) Save(userId uint, content, sourceType, extra string) {
	req := Request{
		UserID:     userId,
		Content:    content,
		SourceType: sourceType,
		Extra:      extra,
		TimeStamp:  time.Now(),
	}
	err := r.DB.Create(&req).Error
	if err != nil {
		s, _ := json.Marshal(req)
		fmt.Println("Failed to save Request[%s]!", string(s))
	}
}
