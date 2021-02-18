package main

import (
	"log"
	"net/http"
	"os"

	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	res := &resolvers.Resolver{
		DatabaseUtil:               database.MongoDatabaseUtil,
		HashUtil:                   hash.BcryptHashUtil,
		IdentifierUtil:             identifier.UUIDIdentifierUtil,
		MailUtil:                   mail.SMTPMailUtil,
		MatchUtil:                  match.RegexpMatchUtil,
		SerializingUtil:            serializing.JSONSerializingUtil,
		TokenUtil:                  token.RngTokenUtil,
		MaxAffinityNumber:          int64(5),
		SecondsLimitAvailability:   int64(2419200),
		SecondsMinimumAvailability: int64(1800),
		SecondsToCooldownReset:     int64(86400),
		SecondsToExpire:            int64(1800),
		SecondsToExpireReset:       int64(86400),
	}

	router := graph.CreateServer(res)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
