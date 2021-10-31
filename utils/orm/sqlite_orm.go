package orm

import (
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
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
		migrateErr := db.AutoMigrate(
			&mails_models.TransientMailMessage{},
			&users_models.Authentication{},
			&users_models.ResetPassword{},
			&users_models.User{},
		)
		if migrateErr != nil {
			panic(migrateErr)
		}
	}

	return err
}

func (p *SqliteOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
