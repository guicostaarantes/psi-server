package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/logging"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

func main() {
	port := os.Getenv("PSI_APP_PORT")
	mongoUri := os.Getenv("PSI_MONGO_URI")
	smtpHost := os.Getenv("PSI_SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("PSI_SMTP_PORT"))
	smtpUser := os.Getenv("PSI_SMTP_USERNAME")
	smtpPass := os.Getenv("PSI_SMTP_PASSWORD")

	loggingUtil := logging.PrintLoggingUtil{}

	databaseUtil := database.MongoDatabaseUtil{
		Context:     context.Background(),
		LoggingUtil: loggingUtil,
	}

	databaseUtil.Connect(mongoUri)

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

	serializingUtil := serializing.JsonSerializingUtil{
		LoggingUtil: loggingUtil,
	}

	tokenUtil := token.RngTokenUtil{
		Runes: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		Size:  64,
	}

	res := &resolvers.Resolver{
		DatabaseUtil:                 &databaseUtil,
		HashUtil:                     hashUtil,
		IdentifierUtil:               identifierUtil,
		MailUtil:                     mailUtil,
		MatchUtil:                    matchUtil,
		SerializingUtil:              serializingUtil,
		TokenUtil:                    tokenUtil,
		MaxAffinityNumber:            int64(5),
		SecondsToCooldownReset:       int64(86400),
		SecondsToExpire:              int64(28800),
		SecondsToExpireReset:         int64(86400),
		TopAffinitiesCooldownSeconds: int64(86400),
	}

	router := graph.CreateServer(res)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
