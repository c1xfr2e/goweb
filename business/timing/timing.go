package timing

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
)

const (
	PeriodDate    = "date"
	PeriodMonth   = "month"
	PeriodQuarter = "quarter"
	PeriodYear    = "year"
)

const (
	DefaultPeriodRangeLength = 90
)

// GetDateRangeOfTable get the min and max date from table in db
func GetDateRangeOfTable(table string, period string, db *gorm.DB,
	filters func(b sq.SelectBuilder) sq.SelectBuilder) (time.Time, time.Time, error) {
	table = fmt.Sprintf(`"%s"`, table)
	sb := sq.Select("min(date) as min, max(date) as max").From(table).Where(sq.Eq{"period": period})
	if filters != nil {
		sb = filters(sb)
	}
	sql, args, err := sb.ToSql()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer rows.Close()
	var min, max time.Time
	if rows.Next() {
		err := rows.Scan(&min, &max)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return min, max, nil
	} else {
		return time.Time{}, time.Time{}, fmt.Errorf("table %s is empty", table)
	}
}

// AlignPeriodRange adjusts time range by period.
func AlignPeriodRange(beginning time.Time, end time.Time, period string) (time.Time, time.Time) {
	if end.IsZero() {
		end = time.Now().UTC()
	}
	end = EndOfPeriod(end, period)

	if beginning.IsZero() {
		beginning = end
		for i := 1; i < DefaultPeriodRangeLength; i++ {
			beginning = Backward(beginning, period)
		}
	} else {
		beginning = BeginningOfPeriod(beginning, period)
	}

	return beginning, end
}

// Backward returns the previous passed time before a given time by specific period.
func Backward(t time.Time, period string) time.Time {
	switch period {
	case PeriodDate:
		justBeforeThisDay := now.New(t).BeginningOfDay().Add(-time.Second)
		return now.New(justBeforeThisDay).BeginningOfDay()

	case PeriodMonth:
		justBeforeThisMonth := now.New(t).BeginningOfMonth().Add(-time.Second)
		return now.New(justBeforeThisMonth).BeginningOfMonth()

	case PeriodQuarter:
		justBeforeThisQuarter := now.New(t).BeginningOfQuarter().Add(-time.Second)
		return now.New(justBeforeThisQuarter).BeginningOfQuarter()

	case PeriodYear:
		justBeforeThisYear := now.New(t).BeginningOfYear().Add(-time.Second)
		return now.New(justBeforeThisYear).BeginningOfYear()

	default:
		logrus.Errorf("Unknow period: %s", period)
	}
	return t
}

// Forward return the next future time after a given time by specific period.
func Forward(t time.Time, period string) time.Time {
	switch period {
	case PeriodDate:
		justAfterThisDay := now.New(t).EndOfDay().Add(time.Second)
		return now.New(justAfterThisDay).BeginningOfDay()

	case PeriodMonth:
		justAfterThisMonth := now.New(t).EndOfMonth().Add(time.Second)
		return now.New(justAfterThisMonth).BeginningOfMonth()

	case PeriodQuarter:
		justAfterThisQuarter := now.New(t).EndOfQuarter().Add(time.Second)
		return now.New(justAfterThisQuarter).BeginningOfQuarter()

	case PeriodYear:
		justAfterThisYear := now.New(t).EndOfYear().Add(time.Second)
		return now.New(justAfterThisYear).BeginningOfYear()

	default:
		logrus.Errorf("Unknow period: %s", period)
	}
	return t
}

// NextPeriodLevel return next period level after given period.
func NextPeriodLevel(period string) string {
	switch period {
	case PeriodDate:
		return PeriodMonth
	case PeriodMonth:
		return PeriodQuarter
	case PeriodQuarter:
		return PeriodYear
	case PeriodYear:
		return PeriodYear
	default:
		logrus.Errorf("Unknow period: %s", period)
	}
	return PeriodDate
}

// BeginningOfPeriod returns the beginning time of period containing the given time
func BeginningOfPeriod(t time.Time, period string) time.Time {
	switch period {
	case PeriodDate:
		return now.New(t).BeginningOfDay()
	case PeriodMonth:
		return now.New(t).BeginningOfMonth()
	case PeriodQuarter:
		return now.New(t).BeginningOfQuarter()
	case PeriodYear:
		return now.New(t).BeginningOfYear()
	default:
		logrus.Errorf("Unknow period: %s", period)
	}
	return time.Now().UTC()
}

// EndOfPeriod returns the end time of period containing the given time
func EndOfPeriod(t time.Time, period string) time.Time {
	switch period {
	case PeriodDate:
		return now.New(t).EndOfDay()
	case PeriodMonth:
		return now.New(t).EndOfMonth()
	case PeriodQuarter:
		return now.New(t).EndOfQuarter()
	case PeriodYear:
		return now.New(t).EndOfYear()
	default:
		logrus.Errorf("Unknow period: %s", period)
	}
	return time.Now().UTC()
}

// FormatTime formats the given time accordingly
func FormatTime(t time.Time, period string) string {
	switch period {
	default:
		fallthrough
	case PeriodDate:
		return t.Format("2006/01/02")
	case PeriodMonth:
		return t.Format("2006/01")
	case PeriodQuarter:
		return fmt.Sprintf("%s/Q%d", t.Format("2006"), int(t.Month()-1)%3+1)
	case PeriodYear:
		return t.Format("2006")
	}
}
