package store

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

var (
	failedErrorMessage = "Failed to write to connection: %v"
)

type Item struct {
	Value []byte
	Flags uint16
}

type StoreInterface interface {
	Set(parts []string, conn net.Conn, reader *bufio.Reader)
	Get(parts []string, conn net.Conn)
}

type Store struct {
	Data map[string]Item
	mu   sync.RWMutex
}

func NewStore() StoreInterface {
	return &Store{
		Data: make(map[string]Item),
		mu:   sync.RWMutex{},
	}
}

//<command name> <key> <flags> <exptime> <byte count> [noreply]\r\n
// <data block>\r\n
func (s *Store) Set(parts []string, conn net.Conn, reader *bufio.Reader) {
	if len(parts) < 5 {
		log.Printf("Invalid set command: %v", parts)
		return
	}

	key := parts[1]

	f, err := strconv.ParseUint(parts[2], 10, 16)
	if err != nil {
		log.Printf("Invalid flags value: %v", parts[2])
		return
	}
	flags := uint16(f)

	bytesCount, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Printf("Invalid bytes count value: %v", parts[4])
		return
	}

	noreply := len(parts) == 6 && parts[5] == "noreply"

	data := make([]byte, bytesCount)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		log.Printf("Failed to read data block: %v", err)
		return
	}

	_, err = reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read trailing newline: %v", err)
		return
	}

	s.mu.Lock()
	s.Data[key] = Item{Value: data, Flags: flags}
	s.mu.Unlock()

	if !noreply {
		_, err := conn.Write([]byte("STORED\r\n"))
		if err != nil {
			log.Printf(failedErrorMessage, err)
		}
	}
}

func (s *Store) Get(parts []string, conn net.Conn) {
	if len(parts) < 2 {
		_, err := conn.Write([]byte("ERROR\r\n"))
		if err != nil {
			log.Printf(failedErrorMessage, err)
		}
		return
	}

	key := parts[1]
	s.mu.RLock()
	item, exists := s.Data[key]
	s.mu.RUnlock()

	if exists {
		response := "VALUE " + key + " " + strconv.Itoa(int(item.Flags)) + " " + strconv.Itoa(len(item.Value)) + "\r\n"
		_, err := conn.Write([]byte(response))
		if err != nil {
			log.Printf(failedErrorMessage, err)
		}
		_, err = conn.Write(item.Value)
		if err != nil {
			log.Printf(failedErrorMessage, err)
		}
		_, err = conn.Write([]byte("\r\n"))
		if err != nil {
			log.Printf(failedErrorMessage, err)
		}
	}

	_, err := conn.Write([]byte("END\r\n"))
	if err != nil {
		log.Printf(failedErrorMessage, err)
	}
}
