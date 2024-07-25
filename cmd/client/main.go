package main

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/solyakov/prxy/pkg/common"
)

type Options struct {
	Listen      string        `short:"l" long:"listen" env:"PRXY_LISTEN" description:"Local address to listen for incoming connections" default:"localhost:8080"`
	Server      string        `short:"s" long:"server" env:"PRXY_SERVER" description:"Remote proxy server address" default:"localhost:8081"`
	Timeout     time.Duration `short:"t" long:"timeout" env:"PRXY_TIMEOUT" description:"Timeout for proxy connections" default:"60s"`
	Buffer      int           `short:"b" long:"buffer" env:"PRXY_BUFFER" description:"Buffer size for data transfer in bytes" default:"32768"`
	Certificate string        `short:"c" long:"certificate" env:"PRXY_CERTIFICATE" description:"Path to the certificate file" default:"client.crt"`
	Key         string        `short:"k" long:"key" env:"PRXY_KEY" description:"Path to the key file" default:"client.key"`
	CA          string        `short:"a" long:"ca" env:"PRXY_CA" description:"Path to the CA certificate file" default:"ca.crt"`
	tlsConfig   *tls.Config
}

func handleClientRequest(client net.Conn, opts Options) {
	defer client.Close()

	server, err := tls.Dial("tcp", opts.Server, opts.tlsConfig)
	if err != nil {
		log.Printf("Failed to connect to proxy server %s: %v", opts.Server, err)
		return
	}
	defer server.Close()

	log.Printf("Opened connection: %s => %s", client.RemoteAddr(), server.RemoteAddr())

	if err := common.Tunnel(server, client, opts.Timeout, opts.Buffer); err != nil {
		log.Printf("Tunneling failed: %v", err)
	}

	log.Printf("Closed connection: %s => %s", client.RemoteAddr(), server.RemoteAddr())
}

func main() {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		}
		log.Fatalf("Error parsing options: %v", err)
	}

	config, err := common.NewTLSConfig(opts.Certificate, opts.Key, opts.CA)
	if err != nil {
		log.Fatalf("Failed to create TLS config: %v", err)
	}
	opts.tlsConfig = config

	listener, err := net.Listen("tcp", opts.Listen)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", opts.Listen, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", opts.Listen)
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleClientRequest(client, opts)
	}
}
