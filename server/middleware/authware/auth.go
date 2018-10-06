package authware

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/bluecover/lm/business/auth"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	userIDKey = "authware.auth.key.userID"
)

func GetCurrentUserID(c *gin.Context) uint {
	v, ok := c.Get(userIDKey)
	if !ok {
		panic("user id not exists inf gin.Context")
	}
	return v.(uint)
}

func LoggedIn(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := isLoggedIn(c, db); err != nil {
			render.Fail(c, err)
			c.Abort()
			return
		}
	}
}

func isLoggedIn(c *gin.Context, db *gorm.DB) error {
	type ReqIdentifyInfo struct {
		Token           string
		FingerprintHash string `valid:"required"`
	}

	req := new(ReqIdentifyInfo)
	token := c.Query("token")
	fp := c.Query("fp")
	if len(token) > 0 && len(fp) > 0 {
		req.Token = token
		fpFrags := strings.Split(fp, "-")
		for _, f := range fpFrags {
			req.FingerprintHash += f
		}
	} else {
		return errors.ErrUnauthorized
	}

	ok, err := govalidator.ValidateStruct(req)
	if !ok || err != nil {
		return errors.ErrUnauthorized
	}

	if len(req.Token) > 0 { // Authenticate request with session token.
		sessionFingerprint := models.GetFingerprintBySessionToken(db, req.FingerprintHash, req.Token)
		if sessionFingerprint == nil {
			return errors.ErrUnauthorized
		}

		c.Set(userIDKey, sessionFingerprint.UserID)
		auth.UpdateSessionToken(db, req.FingerprintHash, req.Token)

	} else { // Request without session token.
		return errors.ErrUnauthorized
	}
	return nil
}

/*
	A sample fingerprint:

{
	"fingerPrint":{
		"b":"8a204daa01a560551d83cb2cb39cbf63fb0e425f6e1d37f2a5c1266f97050270f54a809b33bd7ace48714d2aa955d5faddf2d38c9ab407f8b07007ac5d93e2d4",
		"c":"53e91082c14aaa63953249460db765ca8c52ee386d205543948d2d814958400d2ce0db80ab832e8cc705e03e622d1ce48e71830c28e07e6c265607947b919b69",
		"g":"b1d25916631c6c92780c05e9a06b75b85dc0d06bfa42aeb8da26fcd4eba725c2fa586c77fbe16ffe630ba70b1054d2311d3df65adbd0a7daf92208d81c274045",
		"d":"0f0ac9adef37dd2ce0a64dab2fc58a04043923688373fc1e80c1580c5fdcf8e03a618820f447574a55be43fdd28c657d3f1bb57258a17ecb5bb0028607e7e8c6"
	},

	"info": // tens of thousands long
		"YsqH8YlnPn5N4hfcoYjgEd4zrmJKHa9UnKh4Sl/ydRho9k9RvvarSzgh8L1xkL7US8IZKWhiybVeZaCRm1fj29xYrZGeb1pcWxOcRMBL4Bzdr69QwZx8u7uh799..."
}

*/
