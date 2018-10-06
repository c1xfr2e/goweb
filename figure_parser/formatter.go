package figure_parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	outputTimeFormat = "2006-01-02"
)

var (
	englishPrinter = message.NewPrinter(language.English)
	numberRe       = regexp.MustCompile(`\d+\.?\d*`)
)

func FormatNumber(v interface{}) interface{} {
	switch n := v.(type) {
	case string:
		n = strings.TrimSpace(n)
		if numberRe.MatchString(n) {
			if strings.Contains(n, ".") {
				f, err := strconv.ParseFloat(n, 64)
				if err != nil {
					return v
				}
				return FormatNumber(f)
			} else {
				i, err := strconv.ParseInt(n, 10, 64)
				if err != nil {
					return v
				}
				return FormatNumber(i)
			}
		}
		return v
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
		return englishPrinter.Sprint(n)
	case float32, float64:
		return englishPrinter.Sprintf("%.2f", n)
	default:
		return v
	}
}

func PrintNumber(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprint(FormatNumber(v[0]))
}

func Format(v interface{}) string {
	switch n := v.(type) {
	case time.Time:
		return n.Format(outputTimeFormat)
	default:
		return PrintNumber(v)
	}
}
