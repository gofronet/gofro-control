package main

import (
	"context"
	"errors"
	"gofronet-foundation/gofro-control/internal/nodes/bootstrap"
	"gofronet-foundation/gofro-control/internal/security/certs"
	jwtutils "gofronet-foundation/gofro-control/internal/security/jwt_utils"
	grpcserver "gofronet-foundation/gofro-control/internal/servers/grpc_server"
	httpserver "gofronet-foundation/gofro-control/internal/servers/http_server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	godotenv.Load()

	if err := certs.CreateRootCA(); err != nil {
		panic(err)
	}
	if err := certs.CreateOrEnsureServerCert(); err != nil {
		panic(err)
	}

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errGroup, ctx := errgroup.WithContext(signalCtx)

	jwtSecretManager, err := jwtutils.NewJWTSecretManager()
	if err != nil {
		panic(err)
	}
	inviteStore := bootstrap.NewInviteStore()

	errGroup.Go(func() error {
		err := httpserver.StartHttpServer(ctx, &httpserver.Deps{
			JwtSecretManager: jwtSecretManager,
			InviteStore:      inviteStore,
		})
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server :1337 error: %v", err)
		}
		return err
	})

	errGroup.Go(func() error {
		return grpcserver.StartBootstrapGrpcServer(ctx, &grpcserver.Deps{
			JwtSecretManager: jwtSecretManager,
			InviteStore:      inviteStore,
		})
	})

	if err := errGroup.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		panic(err)
	}

}
