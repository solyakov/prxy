package common

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func Tunnel(end1, end2 net.Conn, timeout time.Duration, size int) error {
	var g errgroup.Group

	g.Go(func() error {
		return copy(end1, end2, timeout, size)
	})
	g.Go(func() error {
		return copy(end2, end1, timeout, size)
	})

	return g.Wait()
}

func copy(source, destination net.Conn, timeout time.Duration, size int) error {
	buf := make([]byte, size)
	for {
		source.SetReadDeadline(time.Now().Add(timeout))
		destination.SetWriteDeadline(time.Now().Add(timeout))
		n, readErr := source.Read(buf)
		if n > 0 {
			written := 0
			for written < n {
				w, writeErr := destination.Write(buf[written:n])
				if writeErr != nil {
					return fmt.Errorf("failed to write data: %v", writeErr)
				}
				written += w
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				return fmt.Errorf("failed to read data: %v", readErr)
			}
			break
		}
	}
	return nil
}

func NewTLSConfig(certificate, key, ca string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCert, err := os.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}, nil
}
