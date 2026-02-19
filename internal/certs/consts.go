package certs

import "gofronet-foundation/gofro-control/internal/constants"

const (
	CertsDir     = constants.AppDataDir + "/certs"
	RootKeyPath  = CertsDir + "/root-ca.key"
	RootCertPath = CertsDir + "/root-ca.crt"
)
