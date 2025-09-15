package logging

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) ToolMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
    return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        start := time.Now()

        m.logger.Printf("Tool call started: tool=%s with following input=%v", req.Params.Name, req.Params.Arguments)

        result, err := next(ctx, req)
        
        duration := time.Since(start)
        if err != nil {
            m.logger.Printf("Tool call failed: tool=%s duration=%v error=%v", req.Params.Name, duration, err)
        } else {
            m.logger.Printf("Tool call completed: tool=%s duration=%v", req.Params.Name, duration)
        }
        
        return result, err
    }
}
