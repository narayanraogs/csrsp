package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	pb "csrsp/server/communication"
	"csrsp/server/global"
	"csrsp/server/rpc"
)

// Version will be set at build time using ldflags.
var Version string

//go:embed all:web
var embeddedFiles embed.FS

func main() {
	listenAddr := flag.String("listen", ":8080", "The address and port to listen on.")
	configPath := flag.String("config", "configuration.json", "Path to the configuration file.")
	flag.Parse()

	log.Printf("Starting server version: %s", Version)

	// Load configuration
	cfg, err := global.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	global.Init(cfg)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCommunicationServer(grpcServer, &rpc.CommunicationServer{})

	// Wrap gRPC server for gRPC-Web support
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	// Get the subtree of the embedded files, so we can serve it from the root.
	webFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(webFS))

	// Create a multiplexer to handle both gRPC-Web and static files
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(r) {
			wrappedGrpc.ServeHTTP(w, r)
			return
		}
		// Fallback to static files
		fileServer.ServeHTTP(w, r)
	})

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	log.Printf("Listening on %s", *listenAddr)
	err = http.ListenAndServe(*listenAddr, c.Handler(handler))
	if err != nil {
		log.Fatal(err)
	}
}
