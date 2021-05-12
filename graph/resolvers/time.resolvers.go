package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"
)

func (r *queryResolver) Time(ctx context.Context) (int64, error) {
	return time.Now().Unix(), nil
}
