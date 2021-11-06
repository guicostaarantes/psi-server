package orm

import (
	"time"

	agreements_models "github.com/guicostaarantes/psi-server/modules/agreements/models"
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
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
				&agreements_models.Agreement{},
				&agreements_models.Term{},
				&appointments_models.Appointment{},
				&characteristics_models.Affinity{},
				&characteristics_models.Characteristic{},
				&characteristics_models.CharacteristicChoice{},
				&characteristics_models.Preference{},
				&cooldowns_models.Cooldown{},
				&mails_models.TransientMailMessage{},
				&profiles_models.Patient{},
				&profiles_models.Psychologist{},
				&treatments_models.Treatment{},
				&treatments_models.TreatmentPriceRange{},
				&treatments_models.TreatmentPriceRangeOffering{},
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
