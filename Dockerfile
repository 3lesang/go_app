FROM golang:1.25.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server cmd/server/main.go

FROM scratch
COPY --from=builder /app/server /
COPY --from=builder /app/docs /

ENTRYPOINT ["/server"]
