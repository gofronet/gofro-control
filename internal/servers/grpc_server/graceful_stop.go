package grpcserver

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func gracefulStopGrpcSerer(ctx context.Context, srv *grpc.Server) {
	<-ctx.Done()
	log.Println("stopping gRPC server")
	srv.GracefulStop()
}
