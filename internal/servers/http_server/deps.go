package httpserver

import (
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap"
	jwtutils "gofronet-foundation/gofro-control/internal/security/jwt_utils"
)

type Deps struct {
	JwtSecretManager *jwtutils.JWTSecretManager
	InviteStore      *bootstrap.InviteStore
}
