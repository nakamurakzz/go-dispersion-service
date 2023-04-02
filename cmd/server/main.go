package main

import (
	"log"

	"github.com/nakamurakzz/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Println("Starting server on :8080")
	log.Fatal(srv.ListenAndServe())
}
