# Memcache Server

**Memcache Server** is a simple in-memory key-value caching server written in Go, inspired by [Memcached](https://memcached.org/).  
It is designed for learning purposes, to understand how caching systems work under the hood, and to practice Go concurrency and networking.

## Features

- Set, Get, and Delete key-value pairs
- In-memory storage
- Expiration (TTL) support
- Concurrent client handling using goroutines
- Simple text-based protocol (similar to Memcached)

## Why This Project?

Building your own caching server helps you understand:

- Networking in Go (TCP/UDP servers)
- Concurrency with goroutines and channels
- Memory management and in-memory data structures
- Basic cache operations like eviction and TTL

## Getting Started

### Requirements
- Go 1.22+ installed

### Run the Server
```bash
git clone https://github.com/dawit-hopes/memcache-server.git
cd memcache-server
go run main.go
