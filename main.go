package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	logger             = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":8080", "gRPC server endpoint") // NOTE: grpc server endpoint options
)

type MoviesServiceServer struct {
	moviesv1.UnimplementedMoviesServiceServer
	mu     *sync.RWMutex
	movies []*moviesv1.Movie
}

func (s *MoviesServiceServer) CreateMovie(ctx context.Context, req *moviesv1.CreateMovieRequest) (*moviesv1.CreateMovieResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	movie := &moviesv1.Movie{
		Id:      uint32(time.Now().Unix()),
		Title:   req.GetTitle(),
		Summary: req.GetSummary(),
		Rating:  req.GetRating(),
	}

	s.movies = append(s.movies, movie)

	res := &moviesv1.CreateMovieResponse{
		Data: movie,
	}

	return res, nil
}

func (s *MoviesServiceServer) ListMovies(ctx context.Context, req *moviesv1.ListMoviesRequest) (*moviesv1.ListMoviesResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := &moviesv1.ListMoviesResponse{
		Data: s.movies,
	}

	return res, nil
}

func (s *MoviesServiceServer) GetMovie(ctx context.Context, req *moviesv1.GetMovieRequest) (*moviesv1.GetMovieResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var movie *moviesv1.Movie

	for _, m := range s.movies {
		if m.Id == req.GetId() {
			movie = m
		}
	}

	if movie == nil {
		return nil, status.Errorf(codes.NotFound, "Movie not found!")
	}

	res := &moviesv1.GetMovieResponse{
		Data: movie,
	}

	return res, nil
}

func (s *MoviesServiceServer) UpdateMovie(ctx context.Context, req *moviesv1.UpdateMovieRequest) (*moviesv1.UpdateMovieResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var movie *moviesv1.Movie
	var movieIdx int

	for i, m := range s.movies {
		if m.Id == req.GetId() {
			movieIdx = i
			movie = m
		}
	}

	if movie == nil {
		return nil, status.Errorf(codes.NotFound, "Movie not found!")
	}

	movie.Title = req.GetTitle()
	movie.Summary = req.GetSummary()
	movie.Rating = req.GetRating()

	s.movies[movieIdx] = movie

	res := &moviesv1.UpdateMovieResponse{
		Data: movie,
	}

	return res, nil

}

func (s *MoviesServiceServer) DeleteMovie(ctx context.Context, req *moviesv1.DeleteMovieRequest) (*moviesv1.DeleteMovieResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	moviesCache := []*moviesv1.Movie{}

	for _, m := range s.movies {
		if m.Id != req.GetId() {
			moviesCache = append(moviesCache, m)
		}
	}

	s.movies = moviesCache

	return &moviesv1.DeleteMovieResponse{}, nil
}

func main() {
	flag.Parse()

	moviesServiceServer := &MoviesServiceServer{
		mu:     &sync.RWMutex{},
		movies: []*moviesv1.Movie{},
	}

	// NOTE: Initialize gRPC Dial Option
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// NOTE: Initialize TCP Connection
	tcp, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		logger.Fatalf("net.Listen: cannot initialize tcp connection")
	}

	// NOTE: Create gRPC Server
	srv := grpc.NewServer()

	// NOTE: Create Mux Handler
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				AllowPartial:    true,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	grpclog.SetLoggerV2(logger)

	// NOTE: Setup context, so the requets can be cancelled
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// NOTE: Run grpc server as go routine
	go func() {
		// NOTE: Register internal servers
		moviesv1.RegisterMoviesServiceServer(srv, moviesServiceServer)

		srv.Serve(tcp)
	}()

	// NOTE: Start HTTP server (and proxy calls to gRPC server endpoint)
	// NOTE: Regsiter request servers
	err = moviesv1.RegisterMoviesServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, dialOptions)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	httpServer.ListenAndServe()
}
