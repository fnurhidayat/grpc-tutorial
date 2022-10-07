package server

import (
	"gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/domain/service"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
)

type MoviesServer struct {
	moviesv1.UnimplementedMoviesServiceServer
	createMovieService service.CreateMovieService
	listMoviesService  service.ListMoviesService
	getMovieService    service.GetMovieService
	updateMovieService service.UpdateMovieService
	deleteMovieService service.DeleteMovieService
}

func NewMoviesServer(
	createMovieService service.CreateMovieService,
	listMoviesService service.ListMoviesService,
	getMovieService service.GetMovieService,
	updateMovieService service.UpdateMovieService,
	deleteMovieService service.DeleteMovieService,
) moviesv1.MoviesServiceServer {
	return &MoviesServer{
		createMovieService: createMovieService,
		listMoviesService:  listMoviesService,
		getMovieService:    getMovieService,
		updateMovieService: updateMovieService,
		deleteMovieService: deleteMovieService,
	}
}
