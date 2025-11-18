package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
)

// Version will be set at build time using ldflags.
var Version string

//go:embed all:web
var embeddedFiles embed.FS

func main() {
	listenAddr := flag.String("listen", ":8080", "The address and port to listen on.")
	flag.Parse()

	log.Printf("Starting server version: %s", Version)

	// Get the subtree of the embedded files, so we can serve it from the root.
	fs, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fs)))

	log.Printf("Listening on %s", *listenAddr)
	err = http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
