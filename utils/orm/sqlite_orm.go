package orm

import (
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteOrmUtil struct {
	dbConn *gorm.DB
}

func (p *SqliteOrmUtil) Connect(dsn string) error {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err == nil {
		p.dbConn = db
		db.AutoMigrate(&users_models.Authentication{})
	}

	return err
}

func (p *SqliteOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
