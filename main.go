package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph/generated"
	"github.com/guicostaarantes/psi-server/graph/generated/model"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
)

const defaultPort = "8082"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := generated.Config{Resolvers: &resolvers.Resolver{}}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []model.Role) (interface{}, error) {
		userID := fmt.Sprintf("%v", ctx.Value("userID"))

		if userID == "" {
			return nil, errors.New("forbidden")
		}

		user, userErr := resolvers.GetUserByIdService.Execute(userID)
		if userErr != nil {
			return nil, errors.New("forbidden")
		}

		for _, v := range role {
			if v == model.Role(user.Role) {
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
				next.ServeHTTP(w, r)
				return
			}

			userID, tokenErr := resolvers.ValidateUserTokenService.Execute(token)
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

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
