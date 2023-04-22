package server

import (
	"context"
	"log"
	"net"
	"regexp"

	"github.com/arianvp/systemd-creds/store"
)

func parsePeerName(s string) (string, string, bool) {
	// Apparently in Go abtract socket names are prefixed with @ instead of 0x00
	matches := regexp.MustCompile("^@.*/unit/(.*)/(.*)$").FindStringSubmatch(s)
	if matches == nil {
		return "", "", false
	}
	unitName := matches[1]
	credID := matches[2]
	return unitName, credID, true
}

func (s *Server) Start(ctx context.Context, connChan <-chan net.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-connChan:
			go s.handleConnection(ctx, conn)
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	unixAddr, ok := conn.RemoteAddr().(*net.UnixAddr)
	if !ok {
		log.Printf("Failed to get peer name: %s", unixAddr.Name)
		return
	}

	unitName, credID, ok := parsePeerName(unixAddr.Name)
	if !ok {
		log.Printf("Failed to parse peer name: %s", unixAddr.Name)
		return
	}
	cred, err := s.Store.Get(ctx, unitName, credID)
	if err != nil {
		log.Printf("Failed to get credential: %v", err)
		return
	}
	if _, err := conn.Write([](byte)(cred)); err != nil {
		log.Printf("Failed to write credential: %v", err)
		return
	}
}

type Server struct {
	Store store.CredStore
}
