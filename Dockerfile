FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o backend main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend .
COPY ./docs ./docs
COPY ./migrations ./migrations
EXPOSE 8080
CMD ["./backend"]
