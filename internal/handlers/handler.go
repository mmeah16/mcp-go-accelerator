package handlers

import (
	"context"
	"mcp-go-accelerator/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
)


func HelloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    arguments := request.GetArguments()
    name, ok := arguments["name"].(string)
    if !ok {
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.TextContent{
                    Type: "text",
                    Text: "Error: name parameter is required and must be a string",
                },
            },
            IsError: true,
        }, nil
    }
 
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Type: "text",
                Text: fmt.Sprintf("Hello to MCP tool, %s!", name),
            },
        },
    }, nil
}

func CalculatorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    num1, ok1 := request.RequireInt("num1")
    num2, ok2 := request.RequireInt("num2")
    operation, ok3 := request.RequireString("operation")

    if ok1 != nil || ok2 != nil || ok3 != nil {
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.TextContent{
                    Type: "text",
                    Text: "Error: Invalid input",
                },
            },
            IsError: true,
        }, nil
    }

    log.Printf("CalculatorHandler: %d %s %d", num1, operation, num2)

    var result int
    switch operation {
    case "add", "+":
        result = num1 + num2
    case "subtract", "-":
        result = num1 - num2
    case "multiply", "*":
        result = num1 * num2
    case "divide", "/":
        if num2 == 0 {
            return &mcp.CallToolResult{
                Content: []mcp.Content{
                    mcp.TextContent{
                        Type: "text",
                        Text: "Error: Division by zero",
                    },
                },
                IsError: true,
            }, nil
        }
        result = num1 / num2
    default:
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.TextContent{
                    Type: "text",
                    Text: "Error: Unknown operation",
                },
            },
            IsError: true,
        }, nil
    }

    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Type: "text",
                Text: fmt.Sprintf("Result: %d", result),
            },
        },
    }, nil
}

func UserHandler(ctx context.Context, request mcp.CallToolRequest, args models.UserInput) (models.UserResponse, error) {
    return models.UserResponse{
        ID:        string(args.FirstName[0]) + "_" + args.LastName + "_" + uuid.New().String(),
        CreatedAt: time.Now(),
    }, nil
}