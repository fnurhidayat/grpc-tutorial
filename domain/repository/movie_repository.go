package repository

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
)

type MovieRepository interface {
	Save(context.Context, *entity.Movie) error
	Get(ctx context.Context, id uint32) (*entity.Movie, error)
	List(context.Context) ([]*entity.Movie, error)
	Destroy(ctx context.Context, id uint32) error
}
