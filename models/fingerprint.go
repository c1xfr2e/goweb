package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Fingerprint represents user auth and access identification flag
type Fingerprint struct {
	ID             uint      `gorm:"primary_key;auto_increment"`
	UserID         uint      `gorm:"column:user_id;not null"`
	HashAll        string    `gorm:"column:hash_all"`
	Origin         string    `gorm:"column:origin;type:text;not null"`
	Token          string    `gorm:"column:token"`
	LastAccessTime time.Time `gorm:"column:last_access_time"`
}

// TableName defines table name
func (Fingerprint) TableName() string {
	return "fingerprints"
}

// GetUserFingerprints fetch fingerprint by user ID
func GetUserFingerprints(db *gorm.DB, userID uint) []Fingerprint {
	var ret []Fingerprint
	db.Where("user_id = ?", userID).Find(&ret)
	return ret
}

// CreateFingerprint creates a new fingerprint and insert into DB for user
func CreateFingerprint(db *gorm.DB, userID uint, fpHash string, origin string) (int, error) {
	r := db.Exec("INSERT INTO fingerprints (user_id, hash_all, origin, last_access_time) "+
		"VALUES(?, ?, ?, ?)", userID, fpHash, origin, time.Now())
	if r.Error != nil {
		return 0, r.Error
	}

	return int(r.RowsAffected), nil
}

// GetFingerprintBySessionToken gets fingerprint by session token
func GetFingerprintBySessionToken(db *gorm.DB, fpHash string, token string) *Fingerprint {
	ret := new(Fingerprint)
	err := db.Where("token=?", token).First(ret).Error
	// ignore hash_all for beta version
	//err := db.Where("hash_all=? AND token=?", fpHash, token).First(ret).Error
	if err != nil {
		return nil
	}
	return ret
}

// ReplaceSessionToken replaces session token and update last access time
func ReplaceSessionToken(db *gorm.DB, userID uint, fingerprint string, token string) error {
	result := db.Model(Fingerprint{}).
		Where("user_id=? AND hash_all=?", userID, fingerprint).
		UpdateColumns(Fingerprint{Token: token, LastAccessTime: time.Now().UTC()})
	return result.Error
}

// UpdateSessionTokenTime updates last access time
func UpdateSessionTokenTime(db *gorm.DB, fpHash string, token string) error {
	err := db.Model(Fingerprint{}).Where("hash_all=? AND token=?", fpHash, token).UpdateColumn(
		"last_access_time", time.Now()).Error
	return err
}

// GetFingerprintOrigin get original data of fingerprint
// TODO: The origin data maybe huge. Query origin alone.
func GetFingerprintOrigin(db *gorm.DB, fpID uint) string {
	return ""
}
