package certs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func CreateOrEnsureServerCert() error {

	dnsListEnv := os.Getenv("TLS_SERVER_DOMAINS")
	ipListEnv := os.Getenv("TLS_SERVER_IPS")

	if ipListEnv == "" && dnsListEnv == "" {
		return errors.New("TLS_SERVER_DNS and TLS_SERVER_IPS env is blank")
	}

	dnsList := strings.Split(dnsListEnv, ",")
	ipListStrings := strings.Split(ipListEnv, ",")

	if len(dnsList) == 1 && dnsList[0] == "" {
		dnsList = nil
	}

	ipList := make([]net.IP, 0, len(ipListStrings))
	for _, ip := range ipListStrings {
		paresedIP := net.ParseIP(ip)
		if paresedIP == nil {
			return fmt.Errorf("invalid IP in TLS_SERVER_IPS: %q", ip)
		}
		ipList = append(ipList, paresedIP)
	}

	if _, err := os.Stat(ServerCertPath); err == nil {
		if _, err2 := os.Stat(ServerKeyPath); err2 == nil {
			return nil
		}
	}

	rootCert, err := readCertPEM(RootCertPath)
	if err != nil {
		return err
	}
	rootKey, err := readECPrivateKeyPKCS8PEM(RootKeyPath)
	if err != nil {
		return err
	}

	// генерим server key
	srvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	serial, err := randSerial()
	if err != nil {
		return err
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.AddDate(3, 0, 0)

	srvTmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "Gofro Panel",
			Organization: []string{"GofroNET"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		IsCA:                  false,
		BasicConstraintsValid: true,

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    dnsList,
		IPAddresses: ipList,
	}

	der, err := x509.CreateCertificate(rand.Reader, srvTmpl, rootCert, srvKey.Public(), rootKey)
	if err != nil {
		return err
	}

	// пишем server.crt
	if err := writePEM(ServerCertPath, "CERTIFICATE", der, 0o644); err != nil {
		return err
	}

	// пишем server.key (PKCS8)
	keyDER, err := x509.MarshalPKCS8PrivateKey(srvKey)
	if err != nil {
		return err
	}
	if err := writePEM(ServerKeyPath, "PRIVATE KEY", keyDER, 0o600); err != nil {
		return err
	}

	return nil
}
