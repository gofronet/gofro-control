package certs

import "gofronet-foundation/gofro-control/internal"

const (
	CertsDir     = internal.AppDataDir + "/certs"
	RootKeyPath  = CertsDir + "/root-ca.key"
	RootCertPath = CertsDir + "/root-ca.crt"
)
