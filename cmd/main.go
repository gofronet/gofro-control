package main

import (
	"gofronet-foundation/gofro-control/internal/certs"
	httpserver "gofronet-foundation/gofro-control/internal/http_server"
)

func main() {
	if err := certs.CreateRootCA(); err != nil {
		panic(err)
	}
	httpserver.StartHttpServer()
}
