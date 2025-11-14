package main

import (
	"flag"

	"github.com/dawit-hopes/memcache-server/internal/server"
	"github.com/dawit-hopes/memcache-server/internal/store"
)

const defaultPort = "11211"

func main() {
	port := flag.String("p", defaultPort, "port to run the server on")
	flag.Parse()

	store := store.NewStore()
	server := server.NewServer(*port, store)
	server.Start()
}
