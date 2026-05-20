# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /server ./server
COPY --from=builder /src/migrations ./migrations
COPY --from=builder /src/docs ./docs

ENV MIGRATIONS_PATH=/app/migrations

EXPOSE 8080

CMD ["./server"]
