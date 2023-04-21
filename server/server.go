package server

import (
	"context"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

type CredStore interface {
	Get(unitName, credID string) (string, error)
}

func parsePeerName(s string) (string, string, bool) {
	matches := regexp.MustCompile("^\x00.*/unit/(.*)/(.*)$").FindStringSubmatch(s)
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
	unitName, credID, ok := parsePeerName(conn.RemoteAddr().String())
	if !ok {
		log.Printf("Failed to parse peer name: %s", conn.RemoteAddr().String())
		return
	}
	cred, err := s.Store.Get(unitName, credID)
	if err != nil {
		log.Printf("Failed to get credential: %v", err)
		return
	}
	if _, err := io.Copy(conn, strings.NewReader(cred)); err != nil {
		log.Printf("Failed to write credential: %v", err)
		return
	}
}

type Server struct {
	Store CredStore
}
