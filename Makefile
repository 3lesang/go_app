.PHONY: dev build run clean swagger

.ONESHELL:
dev:
	export $$(grep -v '^#' .env.local | xargs)
	air

build: clean swagger
	go build -o bin/server
	@ls -lh bin/server | awk '{print $$5}'

run:
	./bin/server

clean:
	rm -rf bin tmp docs