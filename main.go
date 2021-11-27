package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	"github.com/guicostaarantes/psi-server/utils/file_storage"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/logging"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

func main() {
	port := os.Getenv("PSI_APP_PORT")
	smtpHost := os.Getenv("PSI_SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("PSI_SMTP_PORT"))
	smtpUser := os.Getenv("PSI_SMTP_USERNAME")
	smtpPass := os.Getenv("PSI_SMTP_PASSWORD")
	filesBaseFolder := os.Getenv("PSI_FILES_BASE_FOLDER")
	postgresDsn := os.Getenv("PSI_POSTGRES_DSN")

	loggingUtil := logging.PrintLoggingUtil{}

	fileStorageUtil := file_storage.DiskFileStorageUtil{
		BaseFolder:  filesBaseFolder,
		LoggingUtil: loggingUtil,
	}

	hashUtil := hash.BcryptHashUtil{
		Cost:        8,
		LoggingUtil: loggingUtil,
	}

	identifierUtil := identifier.UuidIdentifierUtil{
		LoggingUtil: loggingUtil,
	}

	mailUtil := mail.SmtpMailUtil{
		Host:        smtpHost,
		Port:        smtpPort,
		Username:    smtpUser,
		Password:    smtpPass,
		LoggingUtil: loggingUtil,
	}

	matchUtil := match.RegexpMatchUtil{
		LoggingUtil: loggingUtil,
	}

	ormUtil := orm.PostgresOrmUtil{}

	err := ormUtil.Connect(postgresDsn)
	if err != nil {
		log.Fatalln(err)
	}

	serializingUtil := serializing.JsonSerializingUtil{
		LoggingUtil: loggingUtil,
	}

	tokenUtil := token.RngTokenUtil{
		Runes: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		Size:  64,
	}

	res := &resolvers.Resolver{
		FileStorageUtil:                    fileStorageUtil,
		HashUtil:                           hashUtil,
		IdentifierUtil:                     identifierUtil,
		MailUtil:                           mailUtil,
		MatchUtil:                          matchUtil,
		OrmUtil:                            &ormUtil,
		SerializingUtil:                    serializingUtil,
		TokenUtil:                          tokenUtil,
		MaxAffinityNumber:                  int64(5),
		ScheduleIntervalDuration:           time.Duration(604800) * time.Second,
		ExpireAuthTokenDuration:            time.Duration(28800) * time.Second,
		ExpireResetTokenDuration:           time.Duration(86400) * time.Second,
		InterruptTreatmentCooldownDuration: time.Duration(259200) * time.Second,
		TopAffinitiesCooldownDuration:      time.Duration(86400) * time.Second,
	}

	router := graph.CreateServer(res)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
