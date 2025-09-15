package cli

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)


var (
	mcpClient *client.Client
	ctx       context.Context

	rootCmd = &cobra.Command{
		Use: "mcp-cli",
		Short: "MCP Server CLI",
		Long: "A command-line interface for interacting with the MCP server.",
		Example: "./mcp-cli --help",
		// Inject mcpClient to subsequent commands
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			mcpClient, err = client.NewStreamableHttpClient("http://localhost:8080/mcp")
			if err != nil {
				log.Println("Error creating HTTP client:", err)
				return err
			}
			ctx = context.Background()
			initRequest := mcp.InitializeRequest{}

			if _, err := mcpClient.Initialize(ctx, initRequest); err != nil {
				log.Println("Error initializing MCP client:", err)
			}
			return nil
		},
		// Ensure mcpClient is closed after execution of tools
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if mcpClient != nil {
				mcpClient.Close()
			}
			return nil
		},
	}
)

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Println("Error executing command:", err)
    }
}