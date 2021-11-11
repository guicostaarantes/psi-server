package orm

import "gorm.io/gorm"

type IOrmUtil interface {
	Connect(dsn string) error
	Db() *gorm.DB
}
