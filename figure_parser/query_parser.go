package figure_parser

import (
	"fmt"
	"math"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bluecover/lm/business/timing"
	"github.com/bluecover/lm/util"
	"github.com/jinzhu/gorm"
	gota "github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// Period text const
const (
	PeriodDate    = "date"
	PeriodMonth   = "month"
	PeriodQuarter = "quarter"
	PeriodYear    = "year"

	TimeFormat       = "2006-01-02 15:04:05"
	ResultTimeFormat = "2006/01/02"
)

// QueryParser defines a common data query interface
type QueryParser interface {
	Parse(query map[string]interface{}, args ParseArgs, db *gorm.DB) (map[string]interface{}, error)
}

func newParser(qtype string) (QueryParser, error) {
	switch qtype {
	case "element":
		return elementParser{}, nil
	case "select_column":
		return selectColumnParser{}, nil
	case "aggregate":
		return aggregateParser{}, nil
	case "period_series":
		return periodSeriesParser{}, nil
	case "xox":
		return xoxParser{}, nil
	case "distinct":
		return distinctParser{}, nil
	default:
		return nil, fmt.Errorf("unknow parser type: %s", qtype)
	}
}

func adaptSetFilters(filters []map[string]interface{}) func(sq.SelectBuilder) sq.SelectBuilder {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		return SetFilters(sb, filters)
	}
}

// ParseQuery passes args required to the parser and parse
func ParseQuery(query map[string]interface{}, args ParseArgs, db *gorm.DB) (map[string]interface{}, error) {
	parser, err := newParser(query["type"].(string))
	if err != nil {
		return nil, err
	}

	if p, ok := query["period"]; ok {
		args.Period = p.(string)
	}

	var beginningTime, endTime time.Time

	minTimeOfTable, maxTimeOfTable, err := timing.GetDateRangeOfTable(
		query["table"].(string), args.Period, db, adaptSetFilters(args.Filters))
	if err == nil {
		beginningTime, endTime = timing.AlignPeriodRange(time.Time{}, maxTimeOfTable, args.Period)
	} else {
		logrus.Error("GetDateRangeOfTable error", err)
	}

	if minTimeOfTable.After(beginningTime) {
		beginningTime = timing.BeginningOfPeriod(minTimeOfTable, args.Period)
	}

	if args.Start.IsZero() || args.Start.Before(beginningTime) {
		args.Start = beginningTime
	}
	if args.End.IsZero() || args.End.After(endTime) {
		args.End = endTime
	}

	return parser.Parse(query, args, db)
}

func createPeriodRange(start time.Time, end time.Time, period string) []string {
	prange := make([]string, 0)
	for d := start; d.Before(end) || d == end; d = timing.Forward(d, period) {
		prange = append(prange, timing.FormatTime(d, period))
	}
	return prange
}

func applyArgs(b sq.SelectBuilder, args ParseArgs) sq.SelectBuilder {
	b = b.Where(sq.And{
		sq.Eq{"period": args.Period},
		sq.GtOrEq{"date": args.Start},
		sq.LtOrEq{"date": args.End}},
	)

	b = SetFilters(b, args.Filters)

	return b
}

type distinctParser struct{}

func (distinctParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {

	table := fmt.Sprintf(`"%s"`, query["table"].(string))
	column := query["column"].(string)

	builder := sq.Select("distinct " + column + " as column").From(table)
	statement, sargs, err := builder.ToSql()

	rows, err := db.Raw(statement, sargs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]string, 0)
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			continue
		}
		results = append(results, value)
	}
	return map[string]interface{}{"data": results}, nil
}

type elementParser struct{}

// Parse implement QueryParser.Parse
func (elementParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {

	if one, ok := query["one"].(bool); ok && one {
		args.Start = timing.BeginningOfPeriod(args.End, args.Period)
	}

	table := fmt.Sprintf(`"%s"`, query["table"].(string))
	function := query["function"].(string)

	builder := sq.Select(function).From(table)
	builder = applyArgs(builder, args)
	statement, sargs, err := builder.ToSql()

	rows, err := db.Raw(statement, sargs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v string
	for rows.Next() {
		rows.Scan(&v)
	}
	if len(v) == 0 {
		return map[string]interface{}{"data": nil}, nil
	} else {
		return map[string]interface{}{"data": v}, nil
	}
}

type selectColumnParser struct{}

// Parse implement QueryParser.Parse
func (selectColumnParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {

	table := fmt.Sprintf(`"%s"`, query["table"].(string))
	column := query["column"].(string)

	builder := sq.Select("date", column+" as column").From(table).OrderBy("date ASC")
	builder = applyArgs(builder, args)
	statement, sargs, err := builder.ToSql()

	rows, err := db.Raw(statement, sargs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type R struct {
		Date   string
		Column string
	}
	records := make([]R, 0)
	for rows.Next() {
		var date time.Time
		var column string
		err := rows.Scan(&date, &column)
		if err != nil {
			continue
		}
		records = append(records, R{Date: timing.FormatTime(date, args.Period), Column: column})
	}

	df := gota.LoadStructs(records)
	periodRange := createPeriodRange(args.Start, args.End, args.Period)

	if df.Err != nil {
		resultData := make([]string, len(periodRange))
		for i := range periodRange {
			resultData[i] = "0"
		}
		return map[string]interface{}{
			"date_range": periodRange,
			"data":       resultData,
		}, nil
	}

	dfDateOnly := gota.New(series.New(periodRange, series.String, "Date"))
	dateJoinedColumn := dfDateOnly.LeftJoin(df, "Date").Col("Column").Records()
	replaceNaN(dateJoinedColumn, "-")

	var resultData interface{} = dateJoinedColumn

	_, ok := query["raise_dimension"]
	if ok {
		dim2data := make([][]string, 0)
		dim2data = append(dim2data, dateJoinedColumn)
		resultData = dim2data
	}

	return map[string]interface{}{
		"date_range": periodRange,
		"data":       resultData,
	}, nil
}

type aggregateParser struct{}

// Parse implement QueryParser.Parse
func (aggregateParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {

	table := fmt.Sprintf(`"%s"`, query["table"].(string))
	group := query["group_key"].(string)
	function := query["function"].(string)

	builder := sq.Select("date", group, function).From(table).GroupBy("date", group).OrderBy("date ASC")
	builder = applyArgs(builder, args)
	statement, sargs, err := builder.ToSql()

	rows, err := db.Raw(statement, sargs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type R struct {
		Date  string
		Key   string
		Value string
	}
	groupValueSet := make(map[string]bool)
	records := make([]R, 0)
	for rows.Next() {
		var date time.Time
		var key string
		var value string
		err := rows.Scan(&date, &key, &value)
		if err != nil {
			// TODO: Do not ignore!
			continue
		}
		groupValueSet[key] = true
		records = append(records, R{Date: timing.FormatTime(date, args.Period), Key: key, Value: value})
	}

	df := gota.LoadStructs(records)
	periodRange := createPeriodRange(args.Start, args.End, args.Period)
	dfDateOnly := gota.New(series.New(periodRange, series.String, "Date"))

	groupValueList := make([]string, 0)
	resultData := make([][]string, 0)
	for value := range groupValueSet {
		column := df.Filter(gota.F{Colname: "Key", Comparator: series.Eq, Comparando: value})
		dateJoinedColumn := dfDateOnly.LeftJoin(column, "Date").Col("Value").Records()
		replaceNaN(dateJoinedColumn, "-")
		resultData = append(resultData, dateJoinedColumn)
		groupValueList = append(groupValueList, value)
	}

	return map[string]interface{}{
		"date_range":   periodRange,
		"group_values": groupValueList,
		"data":         resultData,
	}, nil
}

type periodSeriesParser struct{}

// Parse implement QueryParser.Parse
func (periodSeriesParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {

	table := fmt.Sprintf(`"%s"`, query["table"].(string))
	column := query["column"].(string)

	builder := sq.Select("date", column+" as column").From(table).OrderBy("date ASC")
	builder = applyArgs(builder, args)
	statement, sargs, err := builder.ToSql()

	rows, err := db.Raw(statement, sargs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type R struct {
		Date   string
		Column string
	}
	records := make([]R, 0)
	for rows.Next() {
		var date time.Time
		var column string
		err := rows.Scan(&date, &column)
		if err != nil {
			continue
		}
		records = append(records, R{Date: date.Format(ResultTimeFormat), Column: column})
	}

	df := gota.LoadStructs(records)
	periodRange := createPeriodRange(args.Start, args.End, args.Period)
	dfDateOnly := gota.New(series.New(periodRange, series.String, "Date"))
	dfDateJoined := dfDateOnly.LeftJoin(df, "Date")

	arrayPeriodSeries := make([][]string, 0)
	for _, p := range periodRange {
		periodSeries := dfDateJoined.Filter(gota.F{
			Colname: "Date", Comparator: series.Eq, Comparando: p}).Col("Column")
		dateJoinedColumn := periodSeries.Records()
		replaceNaN(dateJoinedColumn, "-")
		arrayPeriodSeries = append(arrayPeriodSeries, dateJoinedColumn)
	}

	return map[string]interface{}{
		"date_range": periodRange,
		"data":       arrayPeriodSeries,
	}, nil
}

type xoxParser struct{}

func xoxName(period string) string {
	switch period {
	case PeriodDate:
		return "Day-on-Day"
	case PeriodMonth:
		return "Month-on-Month"
	case PeriodQuarter:
		return "Quarter-on-Quarter"
	case PeriodYear:
		return "Year-on-Year"
	default:
		return ""
	}
}

// Parse implement QueryParser.Parse
func (p xoxParser) Parse(query map[string]interface{}, args ParseArgs,
	db *gorm.DB) (map[string]interface{}, error) {
	periodLevel := 0
	periodLevelf, ok := query["periodLevel"].(float64)
	if ok {
		periodLevel = int(periodLevelf)
	}

	xPeriod := args.Period
	if periodLevel == 1 {
		xPeriod = timing.NextPeriodLevel(xPeriod)
	}

	periodInQuery, ok := query["period"]
	if ok {
		xPeriod = periodInQuery.(string)
	}

	currentDate := timing.BeginningOfPeriod(args.End, xPeriod)
	args1 := ParseArgs{
		Start:  currentDate,
		End:    currentDate,
		Period: xPeriod,
	}

	r1, err := elementParser{}.Parse(query, args1, db)
	if err != nil {
		return nil, err
	}
	if r1["data"] == nil {
		return p.naResult(xPeriod), nil
	}
	current := cast.ToFloat64(r1["data"])

	prevDate := timing.Backward(currentDate, xPeriod)
	args2 := ParseArgs{
		Start:  prevDate,
		End:    prevDate,
		Period: xPeriod,
	}
	r2, err := elementParser{}.Parse(query, args2, db)
	if err != nil {
		return nil, err
	}
	if r2["data"] == nil {
		return p.naResult(xPeriod), nil
	}
	prev := cast.ToFloat64(r2["data"])

	var change float64
	if prev == 0.0 || math.IsNaN(prev) {
		change = 0.0
	} else {
		change = (current - prev) / prev * 100
	}

	return map[string]interface{}{
		"inc_or_dec": util.CMPFloat(current, prev),
		"change":     fmt.Sprintf("%.2f%%", change),
		"name":       xoxName(xPeriod),
	}, nil
}

func (xoxParser) naResult(period string) map[string]interface{} {
	return map[string]interface{}{
		"inc_or_dec": 1,
		"change":     "N/A",
		"name":       xoxName(period),
	}
}

func replaceNaN(s []string, r string) {
	for i, v := range s {
		if v == "NaN" {
			s[i] = r
		}
	}
}
