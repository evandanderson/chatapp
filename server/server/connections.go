package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

type ConnectionHandler struct {
	conns map[net.Conn]bool
}

func NewConnectionHandler() ConnectionHandler {
	return ConnectionHandler{
		conns: make(map[net.Conn]bool),
	}
}

func (ch ConnectionHandler) CloseConnection(conn net.Conn) {
	log.Printf("Closing %s connection from %s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	delete(ch.conns, conn)
	err := conn.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Connection from %s closed successfully\n", conn.RemoteAddr().String())
	}
}

func (ch ConnectionHandler) HandleConnection(conn net.Conn) {
	defer ch.CloseConnection(conn)
	log.Printf("Accepted %s connection from %s\n", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	ch.conns[conn] = true
	buffer := bytes.NewBuffer(make([]byte, 4096))

	for {
		n, err := conn.Read(buffer.Bytes())
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection from %s closed (EOF)\n", conn.RemoteAddr().String())
				break
			}
			fmt.Printf("Failed to read from buffer: %s\n", err)
		}
		go ch.Publish(conn.RemoteAddr().String(), buffer.Bytes()[:n])
	}
}

func (ch ConnectionHandler) Publish(senderAddr string, payload []byte) {
	for conn := range ch.conns {
		if conn.RemoteAddr().String() != senderAddr {
			if _, err := conn.Write(payload); err != nil { // TODO: Make this async
				log.Println(err)
			}
		}
	}
	log.Printf("[%s]: %s", senderAddr, payload)
}
