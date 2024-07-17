package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/solyakov/prxy/pkg/common"
)

type Options struct {
	Listen      string        `short:"l" long:"listen" env:"PRXY_LISTEN" description:"Local address to listen for incoming connections" default:"localhost:8081"`
	Timeout     time.Duration `short:"t" long:"timeout" env:"PRXY_TIMEOUT" description:"Timeout for proxy connections" default:"60s"`
	Buffer      int           `short:"b" long:"buffer" env:"PRXY_BUFFER" description:"Buffer size for data transfer in bytes" default:"32768"`
	Certificate string        `short:"c" long:"certificate" env:"PRXY_CERTIFICATE" description:"Path to the certificate file" default:"server.crt"`
	Key         string        `short:"k" long:"key" env:"PRXY_KEY" description:"Path to the key file" default:"server.key"`
	CA          string        `short:"a" long:"ca" env:"PRXY_CA" description:"Path to the CA certificate file" default:"ca.crt"`
}

const (
	HTTP_405_METHOD_NOT_ALLOWED = "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
	HTTP_200_OK                 = "HTTP/1.1 200 OK\r\n\r\n"
)

func handleClientRequest(client net.Conn, opts Options) {
	defer client.Close()

	request, err := http.ReadRequest(bufio.NewReader(client))
	if err != nil {
		log.Printf("Failed to read HTTP request: %v", err)
		return
	}

	if request.Method != http.MethodConnect {
		log.Printf("Unsupported HTTP method: %s", request.Method)
		client.Write([]byte(HTTP_405_METHOD_NOT_ALLOWED))
		return
	}

	d := net.Dialer{Timeout: opts.Timeout}
	host, err := d.Dial("tcp", request.Host)
	if err != nil {
		log.Printf("Failed to connect to host %s: %v", request.Host, err)
		return
	}
	defer host.Close()

	log.Printf("Opened connection: %s => %s", client.RemoteAddr(), host.RemoteAddr())

	client.Write([]byte(HTTP_200_OK))
	if err := common.Tunnel(host, client, opts.Timeout, opts.Buffer); err != nil {
		log.Printf("Tunneling failed: %v", err)
	}

	log.Printf("Closed connection: %s => %s", client.RemoteAddr(), host.RemoteAddr())
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

	listener, err := tls.Listen("tcp", opts.Listen, config)
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
