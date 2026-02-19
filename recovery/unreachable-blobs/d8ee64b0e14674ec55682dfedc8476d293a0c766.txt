package grpcserver

import (
	"context"
	apisecurityv1 "gofronet-foundation/gofro-control/gen/go/api/security/v1"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap"
	"gofronet-foundation/gofro-control/internal/security/certs"
	"gofronet-foundation/gofro-control/internal/servers/grpc_server/interceptors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func StartBootstrapGrpcServer(ctx context.Context, deps *Deps) error {
	certs, err := credentials.NewServerTLSFromFile(
		certs.ServerCertPath,
		certs.ServerKeyPath,
	)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.Creds(certs), grpc.UnaryInterceptor(interceptors.UnaryLogging()))

	bootstrapGrpcService := bootstrap.NewBooststrapGrpcService(deps.InviteStore, deps.JwtSecretManager)

	apisecurityv1.RegisterBootstrapServiceServer(server, bootstrapGrpcService)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		return err
	}

	go gracefulStopGrpcSerer(ctx, server)

	log.Printf("gRPC server listening on %s", lis.Addr().String())
	return server.Serve(lis)
}
