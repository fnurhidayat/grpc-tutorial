package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	moviesv1 "gitlab.com/binar-engineering-platform/backend/playground/grpc-tutorial/proto/movies/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	logger             = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":8080", "gRPC server endpoint") // NOTE: grpc server endpoint options
)

type MoviesServiceServer struct {
	moviesv1.UnimplementedMoviesServiceServer
}

func (s *MoviesServiceServer) ListMovies(ctx context.Context, req *moviesv1.ListMoviesRequest) (res *moviesv1.ListMoviesResponse, err error) {
	res = &moviesv1.ListMoviesResponse{
		Data: []*moviesv1.Movie{
			{
				Id:      1,
				Title:   "Pengabdi Setan",
				Summary: "Joko Anwar",
				Rating:  5,
			},
			{
				Id:      2,
				Title:   "Pengabdi Setan 2",
				Summary: "Joko Anwar",
				Rating:  5,
			},
		},
	}

	return res, nil
}

func main() {
	flag.Parse()

	moviesServiceServer := &MoviesServiceServer{}

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
