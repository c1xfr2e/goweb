package figure_parser

import (
	sq "github.com/Masterminds/squirrel"
)

func SetFilters(sb sq.SelectBuilder, filters []map[string]interface{}) sq.SelectBuilder {
	for _, filter := range filters {
		values, ok := filter["values"].([]interface{})
		if !ok || len(values) < 1 {
			continue
		}
		if len(values) == 1 {
			sb = sb.Where(sq.Eq{filter["key"].(string): values[0]})
		} else {
			or := sq.Or{}
			for _, v := range values {
				or = append(or, sq.Eq{filter["key"].(string): v})
			}
			sb = sb.Where(or)
		}
	}

	return sb
}
