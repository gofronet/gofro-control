package httpserver

import (
	"context"
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

func StartHttpServer(ctx context.Context) error {

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

	server := http.Server{
		Addr:    ":1337",
		Handler: r,
	}

	go gracefulDownServer(ctx, &server)

	log.Printf("Starting server on %s", server.Addr)
	return server.ListenAndServe()
}

func gracefulDownServer(ctx context.Context, srv *http.Server) {
	<-ctx.Done()

	log.Printf("shutting down server with addr: %s", srv.Addr)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
