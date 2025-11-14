package main

import (
	"flag"

	"github.com/dawit-hopes/memcache-server/internal/server"
)

const defaultPort = "11211"

func main() {
	port := flag.String("p", defaultPort, "port to run the server on")
	flag.Parse()
	server := server.NewServer(*port)
	server.Start()
}
