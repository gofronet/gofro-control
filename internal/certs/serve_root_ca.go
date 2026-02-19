package certs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type CertsRouter struct {
}

func NewCertsRouter() *CertsRouter {
	return &CertsRouter{}
}

func (*CertsRouter) Register(r chi.Router) {
	r.Route("/certs", func(r chi.Router) {
		r.Get("/root-ca.crt", ServeRootCA)
	})
}

func ServeRootCA(w http.ResponseWriter, r *http.Request) {
	rootCA, err := getRootCA()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/x-pem-file")
	w.WriteHeader(http.StatusOK)
	w.Write(rootCA)
}

func getRootCA() ([]byte, error) {
	cert, err := os.ReadFile(RootCertPath)
	if err != nil {
		return nil, fmt.Errorf("read cert: %w", err)
	}

	return cert, nil
}
