package orm

import (
	agreements_models "github.com/guicostaarantes/psi-server/modules/agreements/models"
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	translations_models "github.com/guicostaarantes/psi-server/modules/translations/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
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
			&translations_models.Translation{},
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
	}

	return err
}

func (p *SqliteOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
