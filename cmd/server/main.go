package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mcp-go-accelerator/pkg/logging"
	"mcp-go-accelerator/pkg/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize components
	logger := log.New(os.Stdout, "[MCP]", log.LstdFlags)

	// Initialize middleware
	loggingMiddleware := logging.NewLoggingMiddleware(logger)

	// Create server with all features 
	s := server.NewMCPServer("HTTP Server", "1.0.0",
		server.WithToolHandlerMiddleware(loggingMiddleware.ToolMiddleware,
		),
	)

	// Add tools to server
	tools.AddProductionTools(s)

	// Start server with graceful shutdown
    startWithGracefulShutdown(s)

}

func startWithGracefulShutdown(s *server.MCPServer) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	http := server.NewStreamableHTTPServer(s)

	go func() {
		log.Println("Starting server... run the inspector to test tools! 🛠️")
		if err := http.Start(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-sigChan 

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := http.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown server:", err)
	}

	log.Println("Server stopped.")

}