package handler

import (
	"encoding/json"

	"github.com/bluecover/lm/business/auth"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zaoshu/hardcore/logging"
)

func RenewToken(db *gorm.DB) gin.HandlerFunc {

	type ReqFP struct {
		FP ReqFingerprint `json:"fingerPrint" valid:"required"`
	}

	return func(c *gin.Context) {
		reqFp := new(ReqFP)
		err := json.Unmarshal([]byte(c.Query("fingerprint")), reqFp)
		if err != nil {
			logging.FromContext(c).Errorf("RenewToken Unmarshal error: %s", err)
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		fingerprint := reqFp.FP
		oldToken := c.PostForm("token")

		reqFpHash := fingerprint.Hash1 + fingerprint.Hash2 + fingerprint.Hash3 + fingerprint.Hash4
		fpThisToken := models.GetFingerprintBySessionToken(db, reqFpHash, oldToken)
		if fpThisToken == nil {
			render.Fail(c, errors.ErrUnauthorized)
			return
		}

		newToken := auth.NewSessionToken(db, fpThisToken.UserID, reqFpHash)
		if newToken == "" {
			logging.FromContext(c).Error("RenewToken NewSessionToken error")
			render.Fail(c, errors.ErrUnauthorized)
			return
		}

		render.OK(c, gin.H{"token": newToken})
	}
}
