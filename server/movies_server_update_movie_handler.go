package server

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

func (s *MoviesServer) UpdateMovie(ctx context.Context, req *moviesv1.UpdateMovieRequest) (*moviesv1.UpdateMovieResponse, error) {
	result, err := s.updateMovieService.Call(ctx, &service.UpdateMovieParams{
		Id:      req.GetId(),
		Title:   req.GetTitle(),
		Summary: req.GetSummary(),
		Rating:  req.GetRating(),
	})
	if err != nil {
		return nil, err
	}

	res := &moviesv1.UpdateMovieResponse{
		Data: &moviesv1.Movie{
			Id:      result.Id,
			Title:   result.Title,
			Summary: result.Summary,
			Rating:  result.Rating,
		},
	}

	return res, nil
}
