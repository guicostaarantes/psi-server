package orm

import (
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
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
			migrateErr := db.AutoMigrate(
				&appointments_models.Appointment{},
				&mails_models.TransientMailMessage{},
				&profiles_models.Patient{},
				&profiles_models.Psychologist{},
				&treatments_models.Treatment{},
				&users_models.Authentication{},
				&users_models.ResetPassword{},
				&users_models.User{},
			)
			if migrateErr != nil {
				panic(migrateErr)
			}
			break
		}

	}

	return ormErr
}

func (p *PostgresOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
