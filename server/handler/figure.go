package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bluecover/lm/business/timing"
	"github.com/bluecover/lm/figure_parser"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/errors"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/oliveagle/jsonpath"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

const reqTimeFormat = "2006/1/2"

func GetFinger(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		idstr := c.Query("ids")
		if len(idstr) == 0 {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}

		var period string
		var selfDefinedTime = false
		datetype := c.Query("dateType")
		switch datetype {
		case "1":
			period = timing.PeriodDate
		case "2":
			period = timing.PeriodMonth
		case "3":
			period = timing.PeriodQuarter
		case "4":
			selfDefinedTime = true
			period = c.Query("customDateType")
			if period == "day" {
				period = "date"
			}
		}
		beginningTime := time.Time{}
		endTime := time.Time{}
		if selfDefinedTime {
			dates := strings.Split(c.Query("dateRange"), ",")
			if len(dates) != 2 {
				render.Fail(c, errors.ErrInvalidParameters)
				return
			}
			t1, err1 := time.Parse(reqTimeFormat, dates[0])
			t2, err2 := time.Parse(reqTimeFormat, dates[1])
			if err1 != nil || err2 != nil {
				render.Fail(c, errors.ErrInvalidParameters)
				logrus.Errorf("time.Parse error: %s %s", err1, err2)
				return
			}
			beginningTime = t1
			endTime = t2
		}

		parseArgs := figure_parser.ParseArgs{
			Start:  beginningTime,
			End:    endTime,
			Period: period,
		}
		filters := make([]map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(c.Query("filters")), &filters); err == nil {
			parseArgs.Filters = filters
		}

		figureIds := strings.Split(idstr, ",")

		figures := make([]map[string]interface{}, 0)
		for _, id := range figureIds {
			figure := models.GetFigure(db, id)
			if figure == nil {
				continue
			}

			var fj map[string]interface{}
			err := json.Unmarshal([]byte(figure.Data), &fj)
			if err != nil {
				continue
			}

			var parsedFigure map[string]interface{}
			figureType := fj["type"].(string)
			if figureType == "table" {
				page, err := strconv.Atoi(c.Query("page"))
				if err != nil {
					continue
				}
				sortBy := c.Query("sortBy")
				if len(sortBy) == 0 {
					sortBy = "date"
				}

				if !parseArgs.Start.IsZero() && !parseArgs.End.IsZero() {
					parseArgs.Start, parseArgs.End = timing.AlignPeriodRange(
						parseArgs.Start,
						parseArgs.End,
						parseArgs.Period,
					)
				}
				parsedFigure, err = figure_parser.ParseTable(fj, parseArgs, page, 12, sortBy, db)
				if err != nil {
					render.Fail(c, err)
					return
				}
			} else {
				parsedFigure, err = figure_parser.ParseFigureB([]byte(figure.Data), parseArgs, db)
				if err != nil {
					render.Fail(c, err)
					return
				}
			}
			figures = append(figures, parsedFigure)
		}

		render.OK(c, gin.H{"figures": figures})
	}
}

func DataExport(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ts, _ := strconv.ParseInt(c.Query("time"), 10, 64)
		t := time.Unix(ts, 0)
		if t.Add(5 * time.Minute).Before(time.Now()) {
			render.Fail(c, errors.ErrInvalidParameters, true)
			return
		}

		id := c.Query("id")
		if len(id) == 0 {
			render.Fail(c, errors.ErrInvalidParameters, true)
			return
		}

		parseArgs, err := getParseArgsFromRequest(c)
		if err != nil {
			render.Fail(c, err, true)
			return
		}

		figure, err := getFigureFromPageID(db, id)
		if err != nil {
			render.Fail(c, err, true)
			return
		}

		sortBy := c.Query("sortBy")
		if len(sortBy) == 0 {
			sortBy = "date"
		}
		queryResult, err := figure_parser.NewQuery(figure, parseArgs, 0, 0, sortBy).Run(db)
		if err != nil {
			render.Fail(c, err, true)
			return
		}

		data, err := queryResultToXlsx(queryResult)
		if err != nil {
			render.Fail(c, err, true)
			return
		}

		c.Header("Content-Disposition", getContentDisposition(id, parseArgs.Start, parseArgs.End))
		c.Data(http.StatusOK, "application/octet-stream", data)
	}
}

func getFigureFromPageID(db *gorm.DB, id string) (map[string]interface{}, error) {
	r := models.GetFigurePage(db, id)
	if r == nil {
		return nil, errors.ErrInvalidParameters
	}

	var j interface{}
	err := json.Unmarshal([]byte(r.Data), &j)
	if err != nil {
		return nil, err
	}

	id = getTableFigureID(j)
	if len(id) == 0 {
		return nil, errors.ErrInvalidParameters
	}

	figure := models.GetFigure(db, id)
	if figure == nil {
		return nil, errors.ErrInvalidParameters
	}

	f := map[string]interface{}{}
	err = json.Unmarshal([]byte(figure.Data), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getTableFigureID(j interface{}) string {
	v, _ := jsonpath.JsonPathLookup(j, `$.dataView[:].figures[:]`)
	if v == nil || len(v.([]interface{})) == 0 {
		return ""
	}

	for _, j := range v.([]interface{}) {
		for _, value := range j.([]interface{}) {
			m := value.(map[string]interface{})
			if m["type"].(string) == "table" {
				return m["id"].(string)
			}
		}
	}
	return ""
}

func getContentDisposition(figureID string, start, end time.Time) string {
	parts := strings.Split(figureID, ".")
	filename := fmt.Sprintf("%s_%s_%s_%s", parts[0], parts[len(parts)-1], start.Format("02012006"), end.Format("02012006"))
	encodeFilename := url.QueryEscape(filename)
	encodeFilename = strings.Replace(encodeFilename, "+", "%20", -1)
	return fmt.Sprintf("attachment; filename*=UTF-8''%s.%s", encodeFilename, "xlsx")
}

func getParseArgsFromRequest(c *gin.Context) (args figure_parser.ParseArgs, err error) {
	var period string
	var selfDefinedTime = false
	datetype := c.Query("dateType")
	switch datetype {
	case "1":
		period = timing.PeriodDate
	case "2":
		period = timing.PeriodMonth
	case "3":
		period = timing.PeriodQuarter
	case "4":
		selfDefinedTime = true
	}
	endTime := time.Time{}
	beginningTime := time.Time{}
	if selfDefinedTime {
		dates := strings.Split(c.Query("dateRange"), ",")
		if len(dates) != 2 {
			render.Fail(c, errors.ErrInvalidParameters)
			return
		}
		beginningTime, _ = time.Parse(reqTimeFormat, dates[0])
		endTime, _ = time.Parse(reqTimeFormat, dates[1])
	}

	args = figure_parser.ParseArgs{
		Start:  beginningTime,
		End:    endTime,
		Period: period,
	}

	args.Start, args.End = timing.AlignPeriodRange(
		args.Start,
		args.End,
		args.Period,
	)
	return args, nil
}

func queryResultToXlsx(qr figure_parser.QueryResult) ([]byte, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	if len(qr.Columns) == 0 {
		sheet.AddRow().AddCell()
	} else {
		sheet.AddRow().WriteSlice(&qr.Columns, -1)

		for _, row := range qr.Data {
			sheetRow := sheet.AddRow()
			for _, value := range row {
				switch v := value.(type) {
				case time.Time:
					sheetRow.AddCell().SetDateWithOptions(v, xlsx.DateTimeOptions{
						Location:        time.UTC,
						ExcelTimeFormat: "yyyy-mm-dd",
					})
				case float32:
					sheetRow.AddCell().SetFloatWithFormat(float64(v), "0.00")
				case float64:
					sheetRow.AddCell().SetFloatWithFormat(v, "0.00")
				default:
					sheetRow.AddCell().SetValue(v)
				}
			}
		}
	}
	if len(qr.Data) > 0 {
		row := qr.Data[0]
		if len(row) > 0 {
			for i, value := range row {
				switch value.(type) {
				case time.Time:
					sheet.SetColWidth(i, i, 11.0)
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	err = file.Write(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
