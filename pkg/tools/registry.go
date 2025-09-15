package tools

import (
	"mcp-go-accelerator/internal/handlers"

	"github.com/mark3labs/mcp-go/mcp"
)

// Define Tools in Registry
var ToolRegistry = map[string]ToolDefinition {
	"hello_world": {
		Tool: NewHelloWorldTool(),
		Handler: handlers.HelloHandler,
	},
	"calculator": {
		Tool: NewCalculatorTool(),
		Handler: handlers.CalculatorHandler,
	},
	"user_tool": {
		Tool: NewUserTool(),
		Handler: mcp.NewStructuredToolHandler(handlers.UserHandler),
	},
}

