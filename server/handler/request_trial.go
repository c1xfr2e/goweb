//  Created by paincompiler on 26/01/2018

package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/module/crm"
	"github.com/bluecover/lm/module/emailidator"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/microcosm-cc/bluemonday"
	"github.com/zaoshu/hardcore/logging"
)

// RequestTrial handles trial request
func RequestTrial(db *gorm.DB) gin.HandlerFunc {
	type Req struct {
		Name    string `schema:"name" form:"name"`       // 姓名
		Company string `schema:"company" form:"company"` // 公司名称
		Email   string `schema:"email" form:"email"`     // 公司邮箱
	}
	return func(c *gin.Context) {
		p := bluemonday.UGCPolicy()
		email := p.Sanitize(c.PostForm("email"))
		name := p.Sanitize(c.PostForm("name"))
		company := p.Sanitize(c.PostForm("company"))

		if len(name) == 0 || len(company) == 0 || len(email) == 0 {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		if err := validateCommercialEmail(email); err != nil {
			render.Fail(c, err)
			return
		}

		ip := c.ClientIP()
		if len(ip) == 0 {
			fields := strings.Split(c.Request.RemoteAddr, ":")
			ip = fields[0]
			if net.ParseIP(ip) == nil {
				ip = "0.0.0.0"
			}
		}

		err := models.NewTrialRequest(db, email, name, company, ip)
		if err != nil {
			render.Fail(c, err)
			return
		}

		obj := crm.NewLeadObject()
		obj.Record.Name = name
		obj.Record.CompanyName = company
		obj.Record.Email = email
		obj.Record.HighSeaID = crm.GConfig.HighSeaID
		err = obj.Create()
		if err != nil {
			logging.FromContext(c).Errorf("[Trial Request] create lead object failed, %v", err)
			render.Fail(c, ErrCRMNotificationFailed)
			return
		}
		render.OK(c, nil)
	}
}

// Error list
var (
	ErrCRMNotificationFailed = errors.New(http.StatusOK, 3021, "crm record failed")

	ErrInvalidEmailFormat  = errors.New(http.StatusOK, 3011, "invalid email address format")
	ErrInvalidCompanyEmail = errors.New(http.StatusOK, 3012, "invalid company email address")
)

func validateCommercialEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return ErrInvalidEmailFormat
	}
	if emailidator.IsPublic(email) {
		return ErrInvalidCompanyEmail
	}
	return nil
}
