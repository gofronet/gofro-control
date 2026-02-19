package certs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	"os"
	"time"
)

func CreateRootCA() error {

	if ensureRootCA() == nil {
		return nil
	}

	if err := os.MkdirAll(CertsDir, 0o755); err != nil {
		return err
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	serial, err := randSerial()
	if err != nil {
		return err
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.AddDate(10, 0, 0)

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "GofroPanelROOT",
			Organization: []string{"GofroNET"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		IsCA:                  true,
		BasicConstraintsValid: true,

		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign,

		// Allow issuing intermediates later (pathLen=1). If you NEVER want intermediates, set to 0.
		MaxPathLen:     0,
		MaxPathLenZero: false,
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, priv.Public(), priv)
	if err != nil {
		return err
	}

	if err := writePEM(RootCertPath, "CERTIFICATE", der, 0o644); err != nil {
		return err
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	if err := writePEM(RootKeyPath, "PRIVATE KEY", keyDER, 0o600); err != nil {
		return err
	}

	return nil
}

func ensureRootCA() error {
	_, err := os.ReadFile(RootCertPath)
	if err != nil {
		return fmt.Errorf("read cert: %w", err)
	}
	_, err = os.ReadFile(RootKeyPath)
	if err != nil {
		return fmt.Errorf("read key: %w", err)
	}
	return nil
}
