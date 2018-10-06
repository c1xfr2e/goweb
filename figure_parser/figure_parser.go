package figure_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	identifierFlag = '#'
	queryKey       = "#query"
	templateKey    = "template"
)

var funcMap = template.FuncMap{"printNumber": PrintNumber}

func matchQueryFlag(s string) bool {
	return len(s) > 0 && s[0] == identifierFlag
}

type QueryTag struct {
	Container interface{}
	Index     int
	Key       string
	Value     string
}

func findQueryTagsInMap(container map[string]interface{}, queryTags *[]QueryTag) {
	for name, obj := range container {
		switch obj.(type) {
		case map[string]interface{}:
			findQueryTagsInMap(obj.(map[string]interface{}), queryTags)

		case []interface{}:
			findQueryTagsInSlice(obj.([]interface{}), queryTags)

		case string:
			if matchQueryFlag(obj.(string)) {
				*queryTags = append(*queryTags,
					QueryTag{Container: container, Key: name, Value: obj.(string)})
			}
		}
	}
}

func findQueryTagsInSlice(container []interface{}, queryTags *[]QueryTag) {
	for index, obj := range container {
		switch obj.(type) {
		case map[string]interface{}:
			findQueryTagsInMap(obj.(map[string]interface{}), queryTags)

		case []interface{}:
			findQueryTagsInSlice(obj.([]interface{}), queryTags)

		case string:
			if matchQueryFlag(obj.(string)) {
				*queryTags = append(*queryTags,
					QueryTag{Container: container, Index: index, Value: obj.(string)})
			}
		}
	}
}

func replaceQueryTags(queryTags []QueryTag, queryResults map[string]interface{}) {
	for _, tag := range queryTags {
		keys := strings.Split(tag.Value, ".")
		switch tag.Container.(type) {
		case map[string]interface{}:
			container := tag.Container.(map[string]interface{})
			if len(keys) == 2 {
				container[tag.Key] = FormatNumber(queryResults[keys[1]])
			} else if len(keys) == 3 {
				container[tag.Key] = FormatNumber(queryResults[keys[1]].(map[string]interface{})[keys[2]])
			}
		case []interface{}:
			container := tag.Container.([]interface{})
			if len(keys) == 2 {
				container[tag.Index] = FormatNumber(queryResults[keys[1]])
			} else if len(keys) == 3 {
				container[tag.Index] = FormatNumber(queryResults[keys[1]].(map[string]interface{})[keys[2]])
			}
		}
	}
}

// ParseArgs represents common args for figure parser
type ParseArgs struct {
	Start   time.Time
	End     time.Time
	Period  string
	Filters []map[string]interface{}
}

// ParseFigureB
func ParseFigureB(figBytes []byte, args ParseArgs, db *gorm.DB) (map[string]interface{}, error) {
	var root map[string]interface{}
	err := json.Unmarshal(figBytes, &root)
	if err != nil {
		return nil, err
	}

	usingTemplate := false
	tflag, ok := root["template"]
	if ok {
		b, ok := tflag.(bool)
		if ok {
			usingTemplate = b
		}
	}

	queries, ok := root[queryKey].(map[string]interface{})
	if !ok {
		return root, nil
	}
	delete(root, queryKey)

	queryTags := make([]QueryTag, 0)
	findQueryTagsInMap(root, &queryTags)
	if len(queryTags) == 0 && !usingTemplate {
		return root, nil
	}

	queryResults := make(map[string]interface{})
	for k, q := range queries {
		query, ok := q.(map[string]interface{})
		if ok {
			result, err := ParseQuery(query, args, db)
			if err != nil {
				return nil, err
			}
			queryResults[k] = result
		}
	}
	if len(queryResults) == 0 {
		queryResults, err = ParseQuery(queries, args, db)
		if err != nil {
			return nil, err
		}
	}

	figureType, ok := root["type"].(string)
	if ok {
		switch figureType {
		case "PieChart":
			parsePieChart(queryResults)
		case "LadderChart.Abs":
			parseLadderChart(queryResults)
		}
	}

	if usingTemplate {
		figstr := strings.Replace(string(figBytes), `\"`, `"`, -1)
		t := template.Must(template.New(root["id"].(string)).Funcs(funcMap).Parse(figstr))
		var buf bytes.Buffer
		err = t.Execute(&buf, queryResults)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(buf.String()), &root)
		if err != nil {
			return nil, err
		}
		delete(root, queryKey)
		delete(root, templateKey)
	} else {
		replaceQueryTags(queryTags, queryResults)
	}

	return root, nil
}

func parsePieChart(queryResult map[string]interface{}) {
	result := make(map[string]string)
	groupValues := queryResult["group_values"].([]string)
	data := queryResult["data"].([][]string)
	sums := make(map[string]float64)
	total := 0.0
	for k, v := range groupValues {
		vsum := 0.0
		for _, d := range data[k] {
			f, e := strconv.ParseFloat(d, 64)
			if e != nil {
				continue
			}
			if math.IsNaN(f) {
				f = 0.0
			}
			vsum += f
			total += f
		}
		sums[v] = vsum
	}

	for k, v := range sums {
		result[k] = fmt.Sprintf("%.2f", v/total*100)
	}
	queryResult["data"] = result
}

func parseLadderChart(queryResult map[string]interface{}) {
	incr, ok1 := queryResult["increase"].(map[string]interface{})
	decr, ok2 := queryResult["decrease"].(map[string]interface{})
	if !ok1 || !ok2 {
		logrus.Error("incorrect ladder chart query")
		return
	}
	incd, ok1 := incr["data"].([]string)
	decd, ok2 := decr["data"].([]string)
	if !ok1 || !ok2 {
		logrus.Error("incorrect ladder chart query data")
		return
	}
	for i := range incd {
		inc := strToFloat(incd[i])
		dec := strToFloat(decd[i])
		diff := inc - dec
		if diff > 0 {
			incd[i] = fmt.Sprintf("%.2f", diff)
			decd[i] = "-"
		} else if diff < 0 {
			incd[i] = "-"
			decd[i] = fmt.Sprintf("%.2f", -diff)
		} else {
			incd[i] = "-"
			decd[i] = "-"
		}
	}
}

func strToFloat(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return 0.0
	}
	return f
}
