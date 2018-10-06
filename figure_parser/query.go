package figure_parser

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
)

type TableRow []interface{}

type Table []TableRow

type Query struct {
	Table   string
	Columns []string
	Where   []interface{}
	Page    int // start from 1
	Limit   int
	OrderBy []string
}

type QueryResult struct {
	Total   int
	Columns []string
	Data    Table
}

func (q Query) NeedPagination() bool {
	return q.Limit > 0
}

func (q Query) GetOffset() int {
	if q.Page > 1 {
		return q.Limit * (q.Page - 1)
	}
	return 0
}

func (q Query) Run(db *gorm.DB) (qr QueryResult, err error) {
	// get total count of data
	qr.Total, err = q.getTotalCount(db)
	if err != nil {
		return
	}

	sql := sq.Select(q.Columns...).From(quoteString(q.Table))
	for _, w := range q.Where {
		sql = sql.Where(w)
	}
	if q.NeedPagination() {
		sql = sql.Limit(uint64(q.Limit)).Offset(uint64(q.GetOffset()))
	}
	for _, o := range q.OrderBy {
		if len(o) > 0 {
			sql = sql.OrderBy(o)
		}
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return
	}

	rows, err := db.Raw(raw, args...).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	qr.Columns, err = rows.Columns()
	if err != nil {
		return
	}

	holder := make([]interface{}, len(qr.Columns), len(qr.Columns))
	for rows.Next() {
		tableRow := make(TableRow, len(qr.Columns), len(qr.Columns))
		for i := 0; i < len(qr.Columns); i++ {
			holder[i] = &tableRow[i]
		}
		err = rows.Scan(holder...)
		if err != nil {
			return
		}
		qr.Data = append(qr.Data, tableRow)
	}
	return
}

func (q Query) getTotalCount(db *gorm.DB) (total int, err error) {
	sql := sq.Select("count(*) AS cnt").From(quoteString(q.Table))
	for _, w := range q.Where {
		sql = sql.Where(w)
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return
	}

	type R struct {
		Cnt int
	}

	r := R{}
	err = db.Raw(raw, args...).Scan(&r).Error
	total = r.Cnt
	return
}

func quoteString(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}

func NewQuery(figure map[string]interface{}, args ParseArgs, page, limit int, orderBy string) (q Query) {
	q.Table = figure["table"].(string)
	if c, ok := figure["columns"].([]interface{}); ok {
		for _, i := range c {
			if s, ok := i.(string); ok {
				q.Columns = append(q.Columns, s)
			}
		}
	}
	if len(q.Columns) == 0 {
		q.Columns = append(q.Columns, "*")
	}

	q.Where = append(q.Where, sq.Eq{"period": args.Period})
	q.Limit = limit
	q.Page = page
	q.OrderBy = append(q.OrderBy, orderBy)
	return
}
