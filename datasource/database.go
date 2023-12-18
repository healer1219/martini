package datasource

import (
	"github.com/healer1219/martini/global"
	"gorm.io/gorm"
)

func PageDB(page int, pageSize int) *gorm.DB {
	return Page(page, pageSize, global.DB())
}

func Page(page int, pageSize int, db *gorm.DB) *gorm.DB {
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}
