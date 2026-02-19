// TODO: refactor this AI bullshit
// but it works for now

package certs

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"time"
)

type IssueLeafOptions struct {
	NodeID           string
	NodeAddress      string
	Organization     string
	IncludeServerEKU bool
}

// IssueLeafFromCSRDER signs the CSR public key with your Root CA and returns leaf cert in DER.
func IssueLeafFromCSRDER(csrDER []byte, opts IssueLeafOptions) (leafDER []byte, notAfter time.Time, err error) {
	csr, err := ParseCSRFromDER(csrDER)
	if err != nil {
		return nil, time.Time{}, err
	}
	// на всякий случай (если VerifyCSRDer вызывается отдельно — можно убрать)
	if err := csr.CheckSignature(); err != nil {
		return nil, time.Time{}, fmt.Errorf("csr signature invalid: %w", err)
	}
	if len(csr.DNSNames) == 0 && len(csr.IPAddresses) == 0 && len(csr.URIs) == 0 {
		return nil, time.Time{}, fmt.Errorf("csr must contain SAN (DNS/IP/URI)")
	}

	rootCert, err := readCertPEM(RootCertPath)
	if err != nil {
		return nil, time.Time{}, err
	}
	rootKey, err := readECPrivateKeyPKCS8PEM(RootKeyPath)
	if err != nil {
		return nil, time.Time{}, err
	}

	serial, err := randSerial()
	if err != nil {
		return nil, time.Time{}, err
	}

	if opts.Organization == "" {
		opts.Organization = "GofroNET"
	}

	// TODO: change cert lifetime
	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter = notBefore.AddDate(10, 0, 0)

	eku := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	if opts.IncludeServerEKU {
		eku = append(eku, x509.ExtKeyUsageServerAuth)
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   opts.NodeID,
			Organization: []string{opts.Organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		BasicConstraintsValid: true,
		IsCA:                  false,

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: eku,

		IPAddresses: []net.IP{net.ParseIP(opts.NodeAddress)},
		URIs:        csr.URIs,
	}

	leafDER, err = x509.CreateCertificate(rand.Reader, tmpl, rootCert, csr.PublicKey, rootKey)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("create leaf cert: %w", err)
	}

	return leafDER, notAfter, nil
}
