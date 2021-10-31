package orm

import (
	"time"

	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOrmUtil struct {
	dbConn *gorm.DB
}

func (p *PostgresOrmUtil) Connect(dsn string) error {
	var ormErr error

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		ormErr = err

		if err == nil {
			p.dbConn = db
			db.AutoMigrate(&users_models.Authentication{})
			break
		}

	}

	return ormErr
}

func (p *PostgresOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
