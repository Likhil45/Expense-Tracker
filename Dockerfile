# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files and download dependencies (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (assuming main.go is in /app)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/expence-tracker

# Final stage: minimal runtime image
FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/expence-tracker /app/expence-tracker

RUN chmod +x /app/expence-tracker

EXPOSE 50001

ENV TZ=UTC

ENTRYPOINT ["/app/expence-tracker"]
