package certs

import (
	"errors"
)

func VerifyCSRDer(csrDer []byte) error {
	csr, err := ParseCSRFromDER(csrDer)
	if err != nil {
		return err
	}

	if err := csr.CheckSignature(); err != nil {
		return errors.New("csr signature invalid")
	}

	if len(csr.DNSNames) == 0 && len(csr.IPAddresses) == 0 && len(csr.URIs) == 0 {
		return errors.New("csr must contain SAN (DNS/IP/URI)")
	}
	return nil
}
