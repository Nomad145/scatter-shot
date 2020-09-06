package main

import (
	"github.com/michaeljoelphillips/scatter-shot/internal/server"
	"github.com/michaeljoelphillips/scatter-shot/internal/storage"
	"os"
)

func main() {
	port := os.Args[1]
	path := os.Args[2]

	server.NewHttpServer(port, storage.NewStorage(path)).Start()
}
