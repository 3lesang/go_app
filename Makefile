.PHONY: dev build run clean swagger

# Run app with Air (auto reload + auto swagger gen)
dev:
	air

# Build production binary
build: clean swagger
	go build -o bin/app main.go
	upx --best --lzma bin/app
	@echo "Binary size after UPX compression:"
	@ls -lh bin/app | awk '{print $$5}'

# Run compiled binary
run:
	./bin/app

# Clean build artifacts
clean:
	rm -rf bin tmp docs

# Generate swagger docs manually
swagger:
	swag init
