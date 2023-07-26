package datasource

import (
	"gorm.io/gorm"
)

func Page(page int, pageSize int, db *gorm.DB) *gorm.DB {
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}
