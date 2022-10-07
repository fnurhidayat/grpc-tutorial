package service

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type UpdateMovieService interface {
	Call(ctx context.Context, params *UpdateMovieParams) (*UpdateMovieResult, error)
}

type UpdateMovieParams struct {
	Id      uint32
	Title   string
	Summary string
	Rating  uint32
}

// We can aggregate the type
type UpdateMovieResult entity.Movie

type UpdateMovieServiceImpl struct {
	movieRepository repository.MovieRepository
	logger          grpclog.LoggerV2
}

func (s *UpdateMovieServiceImpl) Call(ctx context.Context, params *UpdateMovieParams) (*UpdateMovieResult, error) {
	movie, err := s.movieRepository.Get(ctx, params.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve movie!")
	}

	if movie == nil {
		return nil, status.Errorf(codes.NotFound, "Movie not found!")
	}

	movie.Title = params.Title
	movie.Summary = params.Summary
	movie.Rating = params.Rating

	if err := s.movieRepository.Save(ctx, movie); err != nil {
		s.logger.Errorf("[s.movieRepository.Save] %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to update movie!")
	}

	result := &UpdateMovieResult{
		Id:      movie.Id,
		Title:   movie.Title,
		Summary: movie.Summary,
		Rating:  movie.Rating,
	}

	return result, nil
}

func NewUpdateMovieService(
	movieRepository repository.MovieRepository,
	logger grpclog.LoggerV2,
) UpdateMovieService {
	return &UpdateMovieServiceImpl{
		movieRepository: movieRepository,
		logger:          logger,
	}
}
