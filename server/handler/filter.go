package handler

import (
	"github.com/bluecover/lm/figure_parser"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func Filter(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		figureID := c.Query("id")
		if len(figureID) == 0 {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}
		m := models.GetFigurePage(db, figureID)
		if m == nil {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		figurePage, err := figure_parser.ParseFigureB([]byte(m.Data), figure_parser.ParseArgs{}, db)
		if err != nil {
			logrus.Errorf("parse figure page for filter %s error %s", figureID, err)
			render.Fail(c, errors.ErrInternal)
		}

		filter := figurePage["filter"]
		if filter == nil {
			filter = make([]interface{}, 0)
		}
		render.OK(c, gin.H{"filters": filter})
	}

}
