package cli

import (
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)


var listToolsCmd = &cobra.Command{
		Use: "list-tools",
		Short: "List available tools on MCP Server",
		Long: "This command retrieves and displays a list of all tools available on the MCP server.",
		Example: "./mcp-cli list-tools",
		Run: func(cmd *cobra.Command, args []string) {

			toolsRequest := mcp.ListToolsRequest{}
			tools, err := mcpClient.ListTools(ctx, toolsRequest)

			if err != nil {
				log.Println("Error listing tools:", err)
			}

            log.Printf("Available tools: %d\n", len(tools.Tools))
            for _, tool := range tools.Tools {
                log.Printf("- %s: %s\n", tool.Name, tool.Description)
            }
		},
}

func init() {
	listToolsCmd.Flags().StringP("server", "s", "http://localhost:8080/mcp", "MCP Server URL")
	rootCmd.AddCommand(listToolsCmd)
}