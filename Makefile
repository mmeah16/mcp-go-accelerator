server:
	go run cmd/server/main.go

inspector:
	npx @modelcontextprotocol/inspector go run cmd/server/main.go

build:
	podman build -t mcp-go-accelerator .

run:
	podman run -p 8080:8080 mcp-go-accelerator