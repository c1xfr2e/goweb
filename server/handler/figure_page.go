package handler

import (
	"encoding/json"
	"time"

	"github.com/bluecover/lm/business/timing"
	"github.com/bluecover/lm/figure_parser"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetFingerPage(db *gorm.DB) gin.HandlerFunc {

	type DateRange struct {
		Min string `json:"min"`
		Max string `json:"max"`
	}

	const (
		dateViewFixed = iota
		dateViewDate  = 1 << (iota - 1)
		dateViewMonth
		dateViewQuarter
		dateViewCustom
	)

	var period2int = map[string]int{
		timing.PeriodDate:    dateViewDate,
		timing.PeriodMonth:   dateViewMonth,
		timing.PeriodQuarter: dateViewQuarter,
	}

	return func(c *gin.Context) {
		figureID := c.Query("id")
		if len(figureID) == 0 {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		r := models.GetFigurePage(db, figureID)
		if r == nil {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		var figurePage map[string]interface{}
		err := json.Unmarshal([]byte(r.Data), &figurePage)
		if err != nil {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		var beginning, end time.Time
		var period string
		var dateViewFlag = dateViewCustom
		if table, ok := figurePage["table"].(string); ok {
			periods := []string{timing.PeriodQuarter, timing.PeriodMonth, timing.PeriodDate}
			for _, p := range periods {
				b, e, err := timing.GetDateRangeOfTable(table, p, db, nil)
				if err == nil {
					beginning, end = b, e
					period = p
					dateViewFlag |= period2int[p]
				}
			}
		} else {
			logrus.Errorf("no table in figure page %s", figureID)
		}

		if beginning.IsZero() || end.IsZero() {
			render.Fail(c, errors.ErrNoData)
			return
		}

		dateView, ok := figurePage["dateView"].(float64)
		if ok && dateView != dateViewFixed {
			figurePage["dateView"] = int(dateView) & dateViewFlag
		}
		figurePage["dateRange"] = DateRange{
			Min: beginning.Format(DateFormat),
			Max: end.Format(DateFormat),
		}

		if _, ok := figurePage["#query"]; ok {
			parseArgs := figure_parser.ParseArgs{
				Start:  beginning,
				End:    end,
				Period: period,
			}
			parsedFigurePage, err := figure_parser.ParseFigureB([]byte(r.Data), parseArgs, db)
			if err == nil {
				figurePage = parsedFigurePage
			} else {
				logrus.Errorf("parse figure page %s error %s", figureID, err)
			}
		}

		render.OK(c, figurePage)
	}
}
