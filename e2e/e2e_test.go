package e2e

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/guicostaarantes/psi-server/graph/generated"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
	"github.com/stretchr/testify/require"
)

func TestEnd2End(t *testing.T) {

	res := &resolvers.Resolver{
		DatabaseUtil:           database.MockDatabaseUtil,
		HashUtil:               hash.BcryptHashUtil,
		IdentifierUtil:         identifier.UUIDIdentifierUtil,
		MailUtil:               mail.SMTPMailUtil,
		MatchUtil:              match.RegexpMatchUtil,
		SerializingUtil:        serializing.JSONSerializingUtil,
		TokenUtil:              token.RngTokenUtil,
		SecondsToCooldownReset: int64(86400),
		SecondsToExpire:        int64(1800),
		SecondsToExpireReset:   int64(86400),
	}

	res.CreateUserWithPasswordService().Execute(&users_models.CreateUserWithPasswordInput{
		Email:     "coordinator@psi.com.br",
		Password:  "Abc123!@#",
		FirstName: "Bootstrap",
		LastName:  "User",
		Role:      "COORDINATOR",
	})

	config := generated.Config{Resolvers: res}

	config.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []users_models.Role) (interface{}, error) {
		userID := ctx.Value("userID").(string)

		if userID == "" {
			return nil, errors.New("forbidden")
		}

		user, userErr := res.GetUserByIDService().Execute(userID)
		if userErr != nil {
			return nil, errors.New("forbidden")
		}

		for _, v := range role {
			if v == user.Role {
				return next(ctx)
			}
		}

		return nil, errors.New("forbidden")
	}

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(config)))

	storedVariables := map[string]interface{}{}

	t.Run("should not log in with incorrect email", func(t *testing.T) {
		var response struct {
			AuthenticateUser struct {
				Token     string
				ExpiresAt int64
			}
		}

		err := c.Post(`
		query {
			authenticateUser(input: {
				email: "coordinator@psi.com",
				password: "Abc123!@#",
				ipAddress: "100.100.100.100"
			}) {
				token
				expiresAt
			}
		}
		`, &response)

		require.Equal(t, "", response.AuthenticateUser.Token)
		require.Equal(t, int64(0), response.AuthenticateUser.ExpiresAt)
		require.EqualError(t, err, "[{\"message\":\"incorrect credentials\",\"path\":[\"authenticateUser\"]}]")
	})

	t.Run("should not log in with incorrect password", func(t *testing.T) {
		var response struct {
			AuthenticateUser struct {
				Token     string
				ExpiresAt int64
			}
		}

		err := c.Post(`
		query {
			authenticateUser(input: {
				email: "coordinator@psi.com.br",
				password: "Abc123!@%",
				ipAddress: "100.100.100.100"
			}) {
				token
				expiresAt
			}
		}
		`, &response)

		require.Equal(t, "", response.AuthenticateUser.Token)
		require.Equal(t, int64(0), response.AuthenticateUser.ExpiresAt)
		require.EqualError(t, err, "[{\"message\":\"incorrect credentials\",\"path\":[\"authenticateUser\"]}]")
	})

	t.Run("should log in as bootstrap coordinator", func(t *testing.T) {
		var response struct {
			AuthenticateUser struct {
				Token     string
				ExpiresAt int64
			}
		}

		err := c.Post(`
		query {
			authenticateUser(input: {
				email: "coordinator@psi.com.br",
				password: "Abc123!@#",
				ipAddress: "100.100.100.100"
			}) {
				token
				expiresAt
			}
		}
		`, &response)

		require.NotEqual(t, "", response.AuthenticateUser.Token)
		require.Equal(t, time.Now().Unix()+res.SecondsToExpire, response.AuthenticateUser.ExpiresAt)
		require.Nil(t, err)

		storedVariables["coordinator_token"] = response.AuthenticateUser.Token
	})
}
