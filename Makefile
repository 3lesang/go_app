.PHONY: dev build run clean swagger

# Run app with Air (auto reload + auto swagger gen)
dev:
	go install github.com/air-verse/air@latest
	$(shell go env GOPATH)/bin/air

# Build production binary
build: clean swagger
	go build -o bin/app main.go
	@ls -lh bin/app | awk '{print $$5}'

# Run compiled binary
run:
	./bin/app

# Clean build artifacts
clean:
	rm -rf bin tmp docs

# Generate swagger docs manually
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	$(shell go env GOPATH)/bin/swag init
