package services

import (
	"context"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetUsersByRoleService is a service that gets the all users of a specific role in the database
type GetUsersByRoleService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetUsersByRoleService) Execute(role string) ([]*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	users := []*models.User{}

	cursor, findErr := s.DatabaseUtil.FindMany("users", map[string]interface{}{"role": role})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		user := &models.User{}

		decodeErr := cursor.Decode(user)
		if decodeErr != nil {
			return nil, decodeErr
		}

		users = append(users, user)
	}

	return users, nil

}
