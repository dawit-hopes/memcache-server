package server

import (
	"log"
	"net"
)

type TCPServer struct {
	Port string
}

func NewServer(port string) *TCPServer {
	return &TCPServer{Port: port}
}

func (s *TCPServer) Start() {

	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", s.Port, err)
	}
	defer listener.Close()

	log.Printf("Server is listening on port %s", s.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %s", conn.RemoteAddr().String())
}
