# Build stage
FROM golang:1.25.9 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# Final stage — gunakan golang alpine agar Go compiler tersedia
FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .

EXPOSE 8080

CMD ["./main"]