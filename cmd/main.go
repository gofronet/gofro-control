package main

import (
	"context"
	"errors"
	"gofronet-foundation/gofro-control/internal/certs"
	httpserver "gofronet-foundation/gofro-control/internal/http_server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := certs.CreateRootCA(); err != nil {
		panic(err)
	}

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errGroup, ctx := errgroup.WithContext(signalCtx)

	errGroup.Go(func() error {
		err := httpserver.StartHttpServer(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server :1337 error: %v", err)
		}
		return err
	})

	// errGroup.Go(func() error {
	// 	err := httpserver.StartHttpServer(ctx)
	// 	if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		log.Printf("server :1337 error: %v", err)
	// 	}
	// 	return err
	// })

	if err := errGroup.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		panic(err)
	}

}
