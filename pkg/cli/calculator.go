package cli

import (
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

// Sample implementation of command to interact with calculator tool
var calculatorCmd = &cobra.Command{
		Use: "calculator",
		Short: "Interact with the calculator tool",
		Long: "This command allows you to perform calculations using the calculator tool on the MCP server.",
		Args: cobra.ExactArgs(3),
		Example: `./mcp-cli calculator 5 3 add
				  ./mcp-cli calculator 91 23 subtract
				  ./mcp-cli calculator 40 66 multiply
				  ./mcp-cli calculator 28 7 divide`,
		Run: func(cmd *cobra.Command, args []string) {

			result, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "calculator",
					Arguments: map[string]interface{}{
						"num1": args[0],
						"num2":  args[1],
						"operation":  args[2],
					},
				},
			})

			if err != nil {
				log.Println("Error calling calculator tool:", err)
				return
			}
			
			textContent := result.Content[0].(mcp.TextContent)
			log.Printf("%s", textContent.Text)
		},
}

func init() {
	rootCmd.AddCommand(calculatorCmd)
}