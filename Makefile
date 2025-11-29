.PHONY: dev build run clean swagger

# Run app with Air (auto reload + auto swagger gen)
dev:
	go install github.com/air-verse/air@latest
	$(shell go env GOPATH)/bin/air

# Build production binary
build: clean swagger
	go build -o bin/server
	@ls -lh bin/server | awk '{print $$5}'

# Run compiled binary
run:
	./bin/server

# Clean build artifacts
clean:
	rm -rf bin tmp docs

# Generate swagger docs manually
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	$(shell go env GOPATH)/bin/swag init

sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	$(shell go env GOPATH)/bin/sqlc

goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	$(shell go env GOPATH)/bin/goose
