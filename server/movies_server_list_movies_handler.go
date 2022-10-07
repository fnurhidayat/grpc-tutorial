package server

import (
	"context"

	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

func (s *MoviesServer) ListMovies(ctx context.Context, req *moviesv1.ListMoviesRequest) (*moviesv1.ListMoviesResponse, error) {
	result, err := s.listMoviesService.Call(ctx, &service.ListMoviesParams{})
	if err != nil {
		return nil, err
	}

	res := &moviesv1.ListMoviesResponse{
		Data: []*moviesv1.Movie{},
	}

	for _, m := range result.Movies {
		res.Data = append(res.Data, &moviesv1.Movie{
			Id:      m.Id,
			Title:   m.Title,
			Summary: m.Summary,
			Rating:  m.Rating,
		})
	}

	return res, nil
}
