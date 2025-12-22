package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	pb "csrsp/server/communication"
	"csrsp/server/db"
	"csrsp/server/global"
	"csrsp/server/rpc"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Version will be set at build time using ldflags.
var Version string

//go:embed all:web
var embeddedFiles embed.FS

func setupLogger() {
	logDir := global.App.LogFileDirectory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("server_%s.json", timestamp))

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Create a JSON handler for structured logging
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Redirect standard log output to slog
	log.SetOutput(multiWriter)
	log.SetFlags(0) // Remove standard log flags as slog handles timestamps
}

func main() {
	listenAddr := flag.String("listen", ":8080", "The address and port to listen on.")
	configPath := flag.String("config", "configuration.json", "Path to the configuration file.")
	flag.Parse()

	// Load configuration first to get log directory
	cfg, err := global.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	global.Init(cfg)

	dev, _ := global.LoadDeveloperOptions(cfg.DevOpsPath)
	global.SetDeveloperOptions(*dev)

	// Setup logging immediately after config load
	setupLogger()

	slog.Info("Starting server", "version", Version)

	var ips = make([]string, 0)
	for _, pcc := range global.App.PCCList {
		ips = append(ips, pcc.IPAddress)
	}
	err = db.Init(global.App.DBUser, global.App.DBPassword, global.App.DBName, ips)
	if err != nil {
		slog.Error("Unable to connect to Database", "IPs", ips, "User", global.App.DBUser, "Name", global.App.DBName)
		os.Exit(1)
	}

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

	slog.Info("Listening", "address", *listenAddr)
	err = http.ListenAndServe(*listenAddr, c.Handler(handler))
	if err != nil {
		log.Fatal(err)
	}
}
