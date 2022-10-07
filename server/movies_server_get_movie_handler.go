package server

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

func (s *MoviesServer) GetMovie(ctx context.Context, req *moviesv1.GetMovieRequest) (*moviesv1.GetMovieResponse, error) {
	result, err := s.getMovieService.Call(ctx, &service.GetMovieParams{
		Id: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	res := &moviesv1.GetMovieResponse{
		Data: &moviesv1.Movie{
			Id:      result.Id,
			Title:   result.Title,
			Summary: result.Summary,
			Rating:  result.Rating,
		},
	}

	return res, nil
}
