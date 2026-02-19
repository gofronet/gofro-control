// TODO: refactor this AI bullshit
// but it works for now

package bootstrap

import (
	"context"
	"fmt"
	apisecurityv1 "gofronet-foundation/gofro-control/gen/go/api/security/v1"
	"gofronet-foundation/gofro-control/internal/certs"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap/models"
	jwtutils "gofronet-foundation/gofro-control/internal/security/jwt_utils"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BootstrapGrpcService struct {
	inviteStore      *InviteStore
	jwtSecretManager *jwtutils.JWTSecretManager
	apisecurityv1.UnimplementedBootstrapServiceServer
}

func NewBooststrapGrpcService(inviteStore *InviteStore, jwtSecretManager *jwtutils.JWTSecretManager) *BootstrapGrpcService {
	return &BootstrapGrpcService{
		inviteStore:      inviteStore,
		jwtSecretManager: jwtSecretManager,
	}
}

func (s *BootstrapGrpcService) Bootstrap(ctx context.Context, req *apisecurityv1.BootstrapRequest) (*apisecurityv1.BootstrapResponse, error) {

	if req.BootstrapToken == "" {
		return nil, status.Error(codes.InvalidArgument, "boostrap_token is blank")
	}

	tokenData, err := s.jwtSecretManager.Verify(req.BootstrapToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	inviteID := tokenData["invite_id"].(string)
	invite, err := s.inviteStore.GetInvite(inviteID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if time.Now().After(invite.ExpireIn) {
		return nil, status.Error(codes.PermissionDenied, "invite expired")
	}

	if invite.Status != models.InviteStatusPending {
		return nil, status.Error(codes.PermissionDenied, "invite already used at this step")
	}

	if err := certs.VerifyCSRDer(req.CsrDer); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	nodeID := inviteID
	fmt.Println(invite.NodeAddress)

	leafDER, notAfter, err := certs.IssueLeafFromCSRDER(req.CsrDer, certs.IssueLeafOptions{
		NodeAddress:      invite.NodeAddress,
		NodeID:           nodeID,
		Organization:     "GofroNET",
		IncludeServerEKU: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.inviteStore.DoneInvite(inviteID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apisecurityv1.BootstrapResponse{
		NodeId:      nodeID,
		LeafCertDer: leafDER,
		ExpiresUnix: notAfter.Unix(),
	}, nil

}
