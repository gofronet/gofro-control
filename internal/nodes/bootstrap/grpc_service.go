package bootstrap

import (
	"context"
	apisecurityv1 "gofronet-foundation/gofro-control/gen/go/api/security/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BootstrapGrpcService struct {
	apisecurityv1.UnimplementedBootstrapServiceServer
}

func (s *BootstrapGrpcService) Bootstrap(ctx context.Context, req *apisecurityv1.BootstrapRequest) (*apisecurityv1.BootstrapResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Bootstrap not implemented")
}
