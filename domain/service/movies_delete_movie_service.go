package service

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type DeleteMovieService interface {
	Call(ctx context.Context, params *DeleteMovieParams) error
}

type DeleteMovieParams struct {
	Id uint32
}

type DeleteMovieServiceImpl struct {
	movieRepository repository.MovieRepository
	logger          grpclog.LoggerV2
}

func (s *DeleteMovieServiceImpl) Call(ctx context.Context, params *DeleteMovieParams) error {
	if err := s.movieRepository.Destroy(ctx, params.Id); err != nil {
		s.logger.Errorf("[s.movieRepository.Destroy] %s", err.Error())
		return status.Errorf(codes.Internal, "Failed to destroy movie!")
	}

	return nil
}

func NewDeleteMovieService(
	movieRepository repository.MovieRepository,
	logger grpclog.LoggerV2,
) DeleteMovieService {
	return &DeleteMovieServiceImpl{
		movieRepository: movieRepository,
		logger:          logger,
	}
}
