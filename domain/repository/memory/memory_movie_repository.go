package memory

import (
	"context"
	"sync"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
)

type MemoryMovieRepository struct {
	records []*entity.Movie
	db      *sync.RWMutex
}

// Destroy implements repository.MovieRepository
func (r *MemoryMovieRepository) Destroy(ctx context.Context, id uint32) error {
	// We need to lock this, to prevent race condition
	r.db.Lock()

	// We need to unlock this, after we did write things.
	defer r.db.Unlock()

	// Create new collections
	nmovies := []*entity.Movie{}
	for _, m := range r.records {
		if m.Id != id {
			nmovies = append(nmovies, m)
		}
	}

	// Replace existing collections
	r.records = nmovies

	return nil
}

// Get implements repository.MovieRepository
func (r *MemoryMovieRepository) Get(ctx context.Context, id uint32) (*entity.Movie, error) {
	// We need to lock this, to prevent race condition
	r.db.Lock()

	// We need to unlock this, after we did write things.
	defer r.db.Unlock()

	var movie *entity.Movie
	for _, m := range r.records {
		if m.Id == id {
			movie = m
			break
		}
	}

	return movie, nil
}

// List implements repository.MovieRepository
func (r *MemoryMovieRepository) List(context.Context) ([]*entity.Movie, error) {
	// We need to lock this, to prevent race condition
	r.db.Lock()

	// We need to unlock this, after we did write things.
	defer r.db.Unlock()

	return r.records, nil
}

// Save implements repository.MovieRepository
func (r *MemoryMovieRepository) Save(ctx context.Context, im *entity.Movie) error {
	// We need to lock this, to prevent race condition
	r.db.Lock()

	// We need to unlock this, after we did write things.
	defer r.db.Unlock()

	// Lookup and see if we update an entry
	var isUpdating bool
	for i, m := range r.records {
		if m.Id == im.Id {
			isUpdating = true
			m.Title = im.Title
			m.Summary = im.Summary
			m.Rating = im.Rating
			r.records[i] = m
		}
	}

	// If we don't update it, then save
	if !isUpdating {
		r.records = append(r.records, im)
	}

	return nil
}

func NewMemoryMovieRepository() repository.MovieRepository {
	return &MemoryMovieRepository{
		records: []*entity.Movie{},
		db:      &sync.RWMutex{},
	}
}
