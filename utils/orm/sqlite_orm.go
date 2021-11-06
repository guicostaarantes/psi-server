package orm

import (
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
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
			&appointments_models.Appointment{},
			&characteristics_models.Affinity{},
			&characteristics_models.Characteristic{},
			&characteristics_models.CharacteristicChoice{},
			&characteristics_models.Preference{},
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
	}

	return err
}

func (p *SqliteOrmUtil) Db() *gorm.DB {
	return p.dbConn
}
