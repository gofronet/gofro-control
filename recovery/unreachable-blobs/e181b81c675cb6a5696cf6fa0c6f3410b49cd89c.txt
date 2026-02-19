package bootstrap

import (
	"encoding/json"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap/models"
	jwtutils "gofronet-foundation/gofro-control/internal/security/jwt_utils"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type BootstrapRouter struct {
	jwtSecretManager *jwtutils.JWTSecretManager
	inviteStore      *InviteStore
}

func NewBootstrapRouter(jwtSecretManager *jwtutils.JWTSecretManager, invitesStore *InviteStore) *BootstrapRouter {
	return &BootstrapRouter{
		jwtSecretManager: jwtSecretManager,
		inviteStore:      invitesStore,
	}
}

func (router *BootstrapRouter) Register(r chi.Router) {
	r.Route("/bootstrap", func(r chi.Router) {
		r.Post("/invite", router.InviteNode)
	})
}

func (router *BootstrapRouter) InviteNode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req models.InviteNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	req.NodeAddress = strings.TrimSpace(req.NodeAddress)
	if req.NodeAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "node address is blank"})
		return
	}

	inviteID := uuid.NewString()
	now := time.Now()
	exp := now.Add(time.Hour)

	jwtClaims := jwt.MapClaims{
		"iss":       jwtutils.Issuer,
		"aud":       jwtutils.Audience,
		"iat":       now.Unix(),
		"nbf":       now.Unix(),
		"exp":       exp.Unix(),
		"jti":       inviteID,
		"scope":     "bootstrap",
		"invite_id": inviteID,
	}

	signed, err := router.jwtSecretManager.Sign(jwtClaims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	invite := models.InviteRecord{
		NodeAddress: req.NodeAddress,
		InviteID:    inviteID,
		ExpireIn:    exp,
		Status:      models.InviteStatusPending,
	}

	if err := router.inviteStore.AddInvite(&invite); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"bootstrap_token": signed})
}
