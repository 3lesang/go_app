FROM golang:1.25.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 \
    go build -tags netgo -installsuffix netgo \
    -ldflags="-s -w" \
    -o server cmd/server/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/server /server
COPY --from=builder /app/docs /docs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./server"]
