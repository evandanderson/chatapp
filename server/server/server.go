package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	connHandler ConnectionHandler
}

func NewServer() *Server {
	return &Server{
		connHandler: NewConnectionHandler(),
	}
}

func (s *Server) Listen(network string, port uint) error {
	log.Printf("Creating listener on port %d", port)
	ln, err := net.Listen(network, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("Listening for %s connections at %s\n", network, ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.connHandler.HandleConnection(conn)
	}
}
