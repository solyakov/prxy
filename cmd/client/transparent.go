package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"syscall"
	"unsafe"

	"github.com/solyakov/prxy/pkg/common"
)

const (
	SO_ORIGINAL_DST = 80
)

type sockaddrIn struct {
	family uint16
	port   uint16
	addr   [4]byte
	zero   [8]uint8
}

func getOriginalDestination(conn net.Conn) (string, error) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return "", fmt.Errorf("unsupported connection type")
	}

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return "", fmt.Errorf("failed to get raw connection: %v", err)
	}

	var addr sockaddrIn
	size := uint32(unsafe.Sizeof(addr))
	var sockErr error

	if err := rawConn.Control(func(fd uintptr) {
		_, _, errno := syscall.Syscall6(
			syscall.SYS_GETSOCKOPT,
			fd,
			uintptr(syscall.SOL_IP),
			uintptr(SO_ORIGINAL_DST),
			uintptr(unsafe.Pointer(&addr)),
			uintptr(unsafe.Pointer(&size)),
			0,
		)
		if errno != 0 {
			sockErr = errno
		}
	}); err != nil {
		return "", fmt.Errorf("control failed: %v", err)
	}
	if sockErr != nil {
		return "", fmt.Errorf("failed to get original destination: %v", sockErr)
	}

	ip := net.IPv4(addr.addr[0], addr.addr[1], addr.addr[2], addr.addr[3])
	port := uint16(addr.port>>8) | uint16(addr.port<<8)
	
	return fmt.Sprintf("%s:%d", ip.String(), port), nil
}

func discardHTTPResponseHeaders(conn net.Conn) error {
	const response = "HTTP/1.1 200 OK\r\n\r\n"
	buf := make([]byte, len(response))
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	if string(buf) != response {
		return fmt.Errorf("unexpected response: %s", string(buf))
	}
	return err
}

func handleTransparentRequest(client net.Conn, opts Options) {
	defer client.Close()

	originalDest, err := getOriginalDestination(client)
	if err != nil {
		log.Printf("Failed to get original destination: %v", err)
		return
	}

	server, err := tls.Dial("tcp4", opts.Server, opts.tlsConfig)
	if err != nil {
		log.Printf("Failed to connect to proxy server %s: %v", opts.Server, err)
		return
	}
	defer server.Close()

	connectRequest := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", originalDest, originalDest)
	if _, err := server.Write([]byte(connectRequest)); err != nil {
		log.Printf("Failed to send CONNECT request: %v", err)
		return
	}

	if err := discardHTTPResponseHeaders(server); err != nil {
		log.Printf("Failed to read CONNECT response: %v", err)
		return
	}

	log.Printf("Opened transparent connection: %s => %s => %s", client.RemoteAddr(), server.RemoteAddr(), originalDest)

	if err := common.Tunnel(server, client, opts.Timeout, opts.Buffer); err != nil {
		log.Printf("Tunneling failed: %v", err)
	}

	log.Printf("Closed transparent connection: %s => %s => %s", client.RemoteAddr(), server.RemoteAddr(), originalDest)
}
