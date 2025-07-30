.PHONY: build run clean test

# Build the IRC server
build:
	go build -o techircd

# Run the server
run: build
	./techircd

# Clean build artifacts
clean:
	rm -f techircd

# Test with the test client
test: build
	go run test_client.go

# Build for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o techircd-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o techircd-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o techircd-darwin-amd64

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run with race detection
run-race: build
	go run -race *.go
