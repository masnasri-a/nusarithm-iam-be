FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy go.mod dulu untuk cache
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source
COPY . .

# Build binary
RUN GOOS=linux GOARCH=amd64 go build -o backend main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary & resource
COPY --from=builder /app/backend .
COPY ./docs ./docs
COPY ./migrations ./migrations

EXPOSE 8080
CMD ["./backend"]
