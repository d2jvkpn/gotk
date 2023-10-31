package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func Flip(tx *gorm.DB, pageSize, pageNo int) *gorm.DB {
	if pageNo <= 0 || pageSize <= 0 {
		return tx.Limit(0)
	}

	return tx.Limit(pageSize).Offset((pageNo - 1) * pageSize)
}

func MultiEquals[T comparable](tx *gorm.DB, field string, values []T) *gorm.DB {
	wheres, items := make([]string, 0), make([]any, 0)
	unique := make(map[T]bool, len(values))

	for i := range values {
		if unique[values[i]] {
			continue
		}
		wheres = append(wheres, fmt.Sprintf("%s = ?", field))
		items = append(items, values[i])
		unique[values[i]] = true
	}

	if len(wheres) == 0 {
		tx.Where("0 = 1")
	}

	return tx.Where(strings.Join(wheres, " OR "), items...)
}
