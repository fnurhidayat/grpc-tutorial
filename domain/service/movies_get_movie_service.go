package service

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type GetMovieService interface {
	Call(ctx context.Context, params *GetMovieParams) (*GetMovieResult, error)
}

type GetMovieParams struct {
	Id uint32
}

type GetMovieResult entity.Movie

type GetMovieServiceImpl struct {
	movieRepository repository.MovieRepository
	logger          grpclog.LoggerV2
}

func (s *GetMovieServiceImpl) Call(ctx context.Context, params *GetMovieParams) (*GetMovieResult, error) {
	movie, err := s.movieRepository.Get(ctx, params.Id)
	if err != nil {
		s.logger.Errorf("[s.movieRepository.Get] %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to get movie!")
	}

	result := &GetMovieResult{
		Id:      movie.Id,
		Title:   movie.Title,
		Summary: movie.Summary,
		Rating:  movie.Rating,
	}

	return result, nil
}

func NewGetMovieService(
	movieRepository repository.MovieRepository,
	logger grpclog.LoggerV2,
) GetMovieService {
	return &GetMovieServiceImpl{
		movieRepository: movieRepository,
		logger:          logger,
	}
}
