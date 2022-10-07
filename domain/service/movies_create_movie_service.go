package service

import (
	"context"
	"time"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/entity"
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type CreateMovieService interface {
	Call(ctx context.Context, params *CreateMovieParams) (*CreateMovieResult, error)
}

type CreateMovieParams struct {
	Title   string
	Summary string
	Rating  uint32
}

type CreateMovieResult entity.Movie

type CreateMovieServiceImpl struct {
	movieRepository repository.MovieRepository
	logger          grpclog.LoggerV2
}

func (s *CreateMovieServiceImpl) GetId() uint32 {
	return uint32(time.Now().Unix())
}

func (s *CreateMovieServiceImpl) Call(ctx context.Context, params *CreateMovieParams) (*CreateMovieResult, error) {
	movie := &entity.Movie{
		Id:      s.GetId(),
		Title:   params.Title,
		Summary: params.Summary,
		Rating:  params.Rating,
	}

	if err := s.movieRepository.Save(ctx, movie); err != nil {
		s.logger.Errorf("[s.movieRepository.Save] %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to create movie!")
	}

	result := &CreateMovieResult{
		Id:      movie.Id,
		Title:   movie.Title,
		Summary: movie.Summary,
		Rating:  movie.Rating,
	}

	return result, nil
}

func NewCreateMovieService(
	movieRepository repository.MovieRepository,
	logger grpclog.LoggerV2,
) CreateMovieService {
	return &CreateMovieServiceImpl{
		movieRepository: movieRepository,
		logger:          logger,
	}
}
