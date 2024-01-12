package httpserver

import (
	"net"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.port = net.JoinHostPort("", port)
	}
}
