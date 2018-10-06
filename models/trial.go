//  Created by paincompiler on 28/01/2018

package models

import (
	"net/http"

	"github.com/bluecover/lm/server/errors"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
)

// TrialRequest for trial request record entity
type TrialRequest struct {
	gorm.Model
	Email   string `gorm:"column:email;not null"`
	Name    string `gorm:"column:name;not null"`
	Company string `gorm:"column:company;not null"`
	IP      string `gorm:"column:ip"`
}

// TableName defines table name
func (tr *TrialRequest) TableName() string {
	return "trial_request"
}

// NewTrialRequest insert a request
func NewTrialRequest(db *gorm.DB, email, name, company, ip string) error {
	t := db.Debug().Begin()
	var count int
	err := t.Model(&TrialRequest{}).Where("ip = ?", ip).
		Where("created_at BETWEEN ? AND ?", now.BeginningOfDay(), now.EndOfDay()).
		Count(&count).Error
	if err != nil {
		t.Rollback()
		return err
	}
	if count > 5 {
		t.Rollback()
		return errors.New(http.StatusOK, 4001, "out of limit")
	}
	if !t.Where("lower(email) = ?", email).Find(&User{}).RecordNotFound() {
		t.Rollback()
		return errors.New(http.StatusOK, 4000, "email already registered")
	}
	if !t.Where("lower(email) = ?", email).Find(&TrialRequest{}).RecordNotFound() {
		t.Rollback()
		return errors.New(http.StatusOK, 4002, "email already requested")
	}
	err = t.Create(&TrialRequest{
		Email:   email,
		Name:    name,
		Company: company,
		IP:      ip,
	}).Error
	if err != nil {
		t.Rollback()
		return err
	}
	return t.Commit().Error
}
