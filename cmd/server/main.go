package main

import (
	"log"
	"mcp-go-accelerator/pkg/logging"
	"mcp-go-accelerator/pkg/tools"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Initialize components
	logger := log.New(os.Stdout, "[MCP]", log.LstdFlags)

	// Initialize middleware
	loggingMiddleware := logging.NewLoggingMiddleware(logger)

	s := server.NewMCPServer("HTTP Server", "1.0.0",
			server.WithToolHandlerMiddleware(loggingMiddleware.ToolMiddleware,
		),
	)
    tools.AddProductionTools(s)
    // Create HTTP server with custom routes
    mux := http.NewServeMux()
    
    // Add MCP endpoints
    mux.Handle("/mcp", server.NewStreamableHTTPServer(s))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

    // Add middleware
    handler := addMiddleware(mux)

	log.Println("Starting custom StreamableHTTP server on :8080")
	startServer(s, handler)

}

func startServer(s *server.MCPServer, handler http.Handler) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server... run the inspector to test tools! 🛠️")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigChan 

	log.Println("Shutting down server...")

	log.Println("Server stopped.")
}

func addMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
        next.ServeHTTP(w, r) 
    })
}