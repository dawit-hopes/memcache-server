package server

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/dawit-hopes/memcache-server/internal/store"
)

type TCPServer struct {
	Port  string
	store store.StoreInterface
}

func NewServer(port string, store store.StoreInterface) *TCPServer {
	return &TCPServer{Port: port, store: store}
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
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		command := strings.ToLower(parts[0])

		switch command {
		case "set":
			s.store.Set(parts, conn, reader)
		case "get":
			s.store.Get(parts, conn)
		case "add":
			s.store.Add(parts, conn, reader)
		case "replace":
			s.store.Replace(parts, conn, reader)
		case "delete":
			s.store.Delete(parts, conn)
		default:
			if _, err := conn.Write([]byte("ERROR\r\n")); err != nil {
				log.Printf("Failed to write to connection: %v", err)
				return
			}
		}
	}
}
