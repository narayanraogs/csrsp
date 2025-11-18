package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// Version will be set at build time using ldflags.
var Version string

//go:embed all:web
var embeddedFiles embed.FS

func main() {
	log.Printf("Starting server version: %s", Version)

	// Get the subtree of the embedded files, so we can serve it from the root.
	fs, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fs)))

	log.Println("Listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
