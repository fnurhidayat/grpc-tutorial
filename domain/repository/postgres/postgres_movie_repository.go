package postgres

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/driver"
)

type PostgresMovieRepository struct {
	db driver.DB
}

// Destroy implements repository.MovieRepository
func (*PostgresMovieRepository) Destroy(ctx context.Context, id uint32) error {
	panic("unimplemented")
}

// Get implements repository.MovieRepository
func (*PostgresMovieRepository) Get(ctx context.Context, id uint32) (*entity.Movie, error) {
	panic("unimplemented")
}

// List implements repository.MovieRepository
func (r *PostgresMovieRepository) List(ctx context.Context) ([]*entity.Movie, error) {
	movies := []*entity.Movie{}
	rows, err := r.db.QueryContext(ctx, POSTGRES_MOVIE_REPOSITORY_LIST_SQL)
	if err != nil {
		return movies, err
	}

	for rows.Next() {
		movie := &entity.Movie{}

		if err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Summary,
			&movie.Rating,
		); err != nil {
			return movies, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

// Save implements repository.MovieRepository
func (r *PostgresMovieRepository) Save(ctx context.Context, movie *entity.Movie) error {
	// If it's true, then we must update the data instead of insert
	if movie.Id != 0 {
		_, err := r.db.NamedExecContext(ctx, POSTGRES_MOVIE_REPOSITORY_SAVE_UPDATE_SQL, movie)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := r.db.PrepareNamedContext(ctx, POSTGRES_MOVIE_REPOSITORY_SAVE_INSERT_SQL)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryContext(ctx, movie)
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(&movie.Id); err != nil {
			return err
		}
	}

	return nil
}

func NewPostgresMovieRepository(db driver.DB) repository.MovieRepository {
	return &PostgresMovieRepository{
		db: db,
	}
}
