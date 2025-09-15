package tools

import (
	"context"

	"mcp-go-accelerator/internal/models"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolDefinition struct {
	Tool mcp.Tool
	Handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewHelloWorldTool() mcp.Tool {
    // Define a simple tools
    tool := mcp.NewTool("hello_world",
        mcp.WithDescription("Say hello to someone"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("Name of the person to greet"),
        ),
    )
	return tool 
}

func NewCalculatorTool() mcp.Tool {
    tool := mcp.NewTool("calculator",
        mcp.WithDescription("Perform basic arithmetic operations"),
        mcp.WithNumber("num1",
            mcp.Required(),
            mcp.Description("First number"),
        ),
        mcp.WithNumber("num2",
            mcp.Required(),
            mcp.Description("Second number"),
        ),
        mcp.WithString("operation",
            mcp.Required(),
            mcp.Description("Arithmetic operation to perform"),
        ),
    )
    return tool
}

func NewUserTool() mcp.Tool {
    tool := mcp.NewTool("user_tool",
        mcp.WithDescription("Create new user."),
        mcp.WithInputSchema[models.UserInput](),
    )
    return tool
}

func AddProductionTools(s *server.MCPServer) {
	for _, tool := range ToolRegistry {
		s.AddTool(tool.Tool, tool.Handler)
	}
}