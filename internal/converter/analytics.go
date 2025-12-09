package converter

import (
	"database/sql"
	"fmt"
	"time"
)

func FormatGroupKey(groupKey any, groupBy string) string {
	switch v := groupKey.(type) {
	case time.Time:
		if groupBy == "week" {
			year, week := v.ISOWeek()
			return fmt.Sprintf("%d-W%02d", year, week)
		}
		return v.Format("2006-01-02")
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ToFloatPtr(n sql.NullFloat64) *float64 {
	if !n.Valid {
		return nil
	}
	v := n.Float64
	return &v
}

func ToFloatPtrIf(n sql.NullFloat64, condition bool) *float64 {
	if !n.Valid || !condition {
		return nil
	}
	v := n.Float64
	return &v
}
