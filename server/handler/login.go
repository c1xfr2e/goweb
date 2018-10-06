package handler

import (
	"encoding/json"

	"github.com/bluecover/lm/business/auth"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func Login(db *gorm.DB) gin.HandlerFunc {

	type ReqFP struct {
		FP ReqFingerprint `json:"fingerPrint" valid:"required"`
	}

	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		reqFp := new(ReqFP)
		err := json.Unmarshal([]byte(c.PostForm("fingerprint")), reqFp)
		if err != nil {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		fingerprint := reqFp.FP

		user := auth.GetAuthorizedUser(db, email, password)
		if user == nil {
			render.Fail(c, errors.ErrInvalidNameOrPassword)
			return
		}

		loginFpHash := fingerprint.Hash1 + fingerprint.Hash2 + fingerprint.Hash3 + fingerprint.Hash4

		// TODO: Start a goroutine to fix race condition: check existing and create new.
		fingerprintExists := auth.CheckFingerprintExists(db, user.ID, loginFpHash)
		if !fingerprintExists {
			auth.UpdateUserFingerprint(db, user.ID, loginFpHash, fingerprint.Plain)
		}

		newSessionToken := auth.NewSessionToken(db, user.ID, loginFpHash)
		if newSessionToken == "" {
			logrus.Errorf("failed updating session token for user %d", email)
			render.Fail(c, errors.ErrInternal)
			return
		}

		render.OK(c, gin.H{"token": newSessionToken})
	}
}
