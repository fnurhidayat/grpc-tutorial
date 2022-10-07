package server

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

func (s *MoviesServer) CreateMovie(ctx context.Context, req *moviesv1.CreateMovieRequest) (*moviesv1.CreateMovieResponse, error) {
	result, err := s.createMovieService.Call(ctx, &service.CreateMovieParams{
		Title:   req.GetTitle(),
		Summary: req.GetSummary(),
		Rating:  req.GetRating(),
	})
	if err != nil {
		return nil, err
	}

	res := &moviesv1.CreateMovieResponse{
		Data: &moviesv1.Movie{
			Id:      result.Id,
			Title:   result.Title,
			Summary: result.Summary,
			Rating:  result.Rating,
		},
	}

	return res, nil
}
