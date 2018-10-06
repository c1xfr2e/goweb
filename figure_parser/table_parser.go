package figure_parser

import (
	"github.com/jinzhu/gorm"
)

func ParseTable(fig map[string]interface{}, args ParseArgs, page, limit int, sortby string, db *gorm.DB) (map[string]interface{}, error) {
	queryResult, err := NewQuery(fig, args, page, limit, sortby).Run(db)
	if err != nil {
		return nil, err
	}

	columnData := make([][]string, len(queryResult.Columns))
	for _, row := range queryResult.Data {
		for i, value := range row {
			columnData[i] = append(columnData[i], Format(value))
		}
	}

	figDataColumns := make([]map[string]interface{}, len(queryResult.Columns))
	for i, col := range queryResult.Columns {
		figDataColumns[i] = map[string]interface{}{
			"title":      col,
			"sortSymbol": col,
			"size":       "medium",
			"data":       columnData[i],
		}
	}

	figData, ok := fig["data"].(map[string]interface{})
	if !ok {
		figData = map[string]interface{}{}
	}
	figData["total"] = queryResult.Total
	figData["currentPage"] = page
	figData["columns"] = figDataColumns
	fig["data"] = figData
	return fig, nil
}
