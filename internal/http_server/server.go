package httpserver

import (
	"gofronet-foundation/gofro-control/internal/certs"
	jwtutils "gofronet-foundation/gofro-control/internal/jwt_utils"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func StartHttpServer() error {

	r := chi.NewRouter()
	{
		r.Use(chi_middlewares.RequestID)
		r.Use(chi_middlewares.RealIP)
		r.Use(chi_middlewares.Timeout(60 * time.Second))
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		r.Use(chi_middlewares.Logger)
		r.Use(chi_middlewares.Recoverer)
	}

	inviteStore := bootstrap.NewInviteStore()
	jwtSecretManager, err := jwtutils.NewJWTSecretManager()
	if err != nil {
		panic(err)
	}

	bootstrapRouter := bootstrap.NewBootstrapRouter(jwtSecretManager, inviteStore)
	serveRootCaRouter := certs.NewCertsRouter()
	r.Route("/v1", func(r chi.Router) {
		serveRootCaRouter.Register(r)
		bootstrapRouter.Register(r)
	})

	log.Println("Starting server on :1337")
	log.Fatalln(http.ListenAndServe(":1337", r))
	return nil
}
