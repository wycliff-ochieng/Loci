FROM golang:1.24-bookworm AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#copy source code
COPY  . .

RUN go build -o /dist/main ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /dist/main ./main

COPY --from=builder /app/internal/store/migrations ./internal/store/migrations

EXPOSE 3000

ENTRYPOINT ["/app/main"]
