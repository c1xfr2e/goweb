package handler

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/middleware/authware"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zaoshu/backgroundservice/client"
	"github.com/zaoshu/hardcore/logging"
)

type RequestOp struct {
	DB *gorm.DB
}

var requestEmailFormat = `
Email: %s <br>
Organization: %s <br>
RequestTime: %s <br>
SourceType: %s <br>
Extra: %s <br>
RequestDescription: %s <br>
`

func (r *RequestOp) NewRequest() gin.HandlerFunc {
	type ReqData struct {
		Content string `json:"content" valid:"required"`
		Type    string `json:"type" valid:"required"`
		Extra   string `json:"extra" valid:"required"`
	}

	return func(c *gin.Context) {
		userId := authware.GetCurrentUserID(c)
		sourceType := c.PostForm("type")
		extra := c.PostForm("extra")
		content := c.PostForm("content")

		op := &models.RequestDBOp{r.DB}
		op.Save(userId, content, sourceType, extra)
		user := models.GetUserByID(r.DB, userId)
		if user == nil {
			panic(fmt.Sprintf("Should not be here. User[%d] doesn't exist???", userId))
		}
		requestEmailAddress := os.Getenv("RequestEmailAddress")
		if requestEmailAddress == "" {
			panic("Env RequestEmailAddress shouldn't be empty")
		}
		err := SendRequestEmail(
			fmt.Sprintf(requestEmailFormat,
				user.Email,
				"",
				time.Now().Format(time.RFC850),
				sourceType,
				extra,
				content),
			requestEmailAddress)
		if err != nil {
			logging.FromContext(c).Errorf("Failed to Send Request Email, err[%v]", err)
		}

		render.OK(c, nil)
	}
}

func SendRequestEmail(emailContent, emailAddress string) error {
	return client.SendEmail(context.Background(), "New Customer Request", emailContent, []string{emailAddress}, nil, "", nil)

}
