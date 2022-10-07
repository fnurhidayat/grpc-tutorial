package service

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type ListMoviesService interface {
	Call(ctx context.Context, params *ListMoviesParams) (*ListMoviesResult, error)
}

type ListMoviesParams struct{}

type ListMoviesResult struct {
	Movies []*entity.Movie
}

type ListMoviesServiceImpl struct {
	movieRepository repository.MovieRepository
	logger          grpclog.LoggerV2
}

func (s *ListMoviesServiceImpl) Call(ctx context.Context, params *ListMoviesParams) (*ListMoviesResult, error) {
	movies, err := s.movieRepository.List(ctx)
	if err != nil {
		s.logger.Errorf("[s.movieRepository.List] %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to list movies!")
	}

	result := &ListMoviesResult{
		Movies: movies,
	}

	return result, nil
}

func NewListMoviesService(
	movieRepository repository.MovieRepository,
	logger grpclog.LoggerV2,
) ListMoviesService {
	return &ListMoviesServiceImpl{
		movieRepository: movieRepository,
		logger:          logger,
	}
}
