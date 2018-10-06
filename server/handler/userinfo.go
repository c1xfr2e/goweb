package handler

import (
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/middleware/authware"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserInfo(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := authware.GetCurrentUserID(c)
		user := models.GetUserByID(db, userID)
		if user == nil {
			render.Fail(c, errors.ErrUnauthorized)
			return
		}

		userInfo := gin.H{
			"accountName": user.Email,
			"email":       user.Email,
			"expiresIn":   "2018/12/31",
			"locale":      "en-US",
		}

		render.OK(c, gin.H{"user": userInfo})
	}
}
