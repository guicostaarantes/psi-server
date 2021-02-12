package graph

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph/generated"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
)

// CreateServer will take the resolver object with dependencies and return a Mux router for the GraphQL application
func CreateServer(res *resolvers.Resolver) *chi.Mux {

	bootstrapUser := os.Getenv("PSI_BOOTSTRAP_USER")
	bootstrap := strings.Split(bootstrapUser, "|")

	if len(bootstrap) == 2 {
		fmt.Println("Bootstrap user identified")
		res.CreateUserWithPasswordService().Execute(&users_models.CreateUserWithPasswordInput{
			Email:     bootstrap[0],
			Password:  bootstrap[1],
			FirstName: "Bootstrap",
			LastName:  "User",
			Role:      "COORDINATOR",
		})
	}

	c := generated.Config{Resolvers: res}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []users_models.Role) (interface{}, error) {
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

	router := chi.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				ctx := context.WithValue(r.Context(), "userID", "")
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			userID, tokenErr := res.ValidateUserTokenService().Execute(token)
			if tokenErr != nil {
				ctx := context.WithValue(r.Context(), "userID", "")
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	srv.Use(extension.FixedComplexityLimit(10))

	router.Handle("/", playground.Handler("GraphQL playground", "/gql"))
	router.Handle("/gql", srv)

	return router

}
