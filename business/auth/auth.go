package auth

import (
	"github.com/jinzhu/gorm"

	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/util"
)

const (
	lenSessionToken = 64
)

// GetAuthorizedUser find user from DB by email and check if password matches
func GetAuthorizedUser(db *gorm.DB, email string, plainPassword string) *models.User {
	user := models.GetUserByEmail(db, email)
	if user == nil {
		return nil
	}
	matched := util.ComparePassword(user.Password, plainPassword)
	if !matched {
		return nil
	}
	return user
}

// CheckFingerprintExists checks if fingerprint exists
func CheckFingerprintExists(db *gorm.DB, userID uint, fpHash string) bool {
	fingerprints := models.GetUserFingerprints(db, userID)
	exists := false
	for _, v := range fingerprints {
		if v.HashAll == fpHash {
			exists = true
			break
		}
	}
	return exists
}

// NewSessionToken generate and update session token
func NewSessionToken(db *gorm.DB, userID uint, fpHash string) string {
	newSessionToken := util.GenerateRandomString(lenSessionToken)
	err := models.ReplaceSessionToken(db, userID, fpHash, newSessionToken)
	if err != nil {
		return ""
	}
	return newSessionToken
}

// UpdateSessionToken updates latest access time of specific hash and token
func UpdateSessionToken(db *gorm.DB, fpHash string, token string) error {
	return models.UpdateSessionTokenTime(db, fpHash, token)
}

// UpdateUserFingerprint updates fingerprint
// There'll be a max number threshold of fingerprints per user
func UpdateUserFingerprint(db *gorm.DB, userID uint, fpHash string, fpOrigin string) (int, error) {
	// TODO: Limit max fingerprints per user.
	// Like this:
	//	insert into fingerprints
	//		select THIS where (select count(*) from fingerprints where user_id=THIS.user_id) < 5"

	nCreated, err := models.CreateFingerprint(db, userID, fpHash, fpOrigin)
	return nCreated, err
}

// GetAuthorizedDatasets get targets datasets
// TODO: access permission mechanism
func GetAuthorizedDatasets(db *gorm.DB, userID uint) []models.Dataset {
	return models.GetAllDatasets(db)
}
