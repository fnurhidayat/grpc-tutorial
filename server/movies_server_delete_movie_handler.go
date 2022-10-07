package server

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

func (s *MoviesServer) DeleteMovie(ctx context.Context, req *moviesv1.DeleteMovieRequest) (*moviesv1.DeleteMovieResponse, error) {
	if err := s.deleteMovieService.Call(ctx, &service.DeleteMovieParams{
		Id: req.GetId(),
	}); err != nil {
		return nil, err
	}

	res := &moviesv1.DeleteMovieResponse{}

	return res, nil
}
