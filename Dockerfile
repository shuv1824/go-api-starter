# -----------------------
# Stage 1: Base for development
# -----------------------
FROM golang:1.22-alpine AS dev

# Install essential tools
RUN apk add --no-cache git bash make curl build-base

# Install Air for live reload
RUN go install github.com/cosmtrek/air@latest

# Set the working directory
WORKDIR /app

# Expose application port (change if needed)
EXPOSE 8080

# Add Air config (optional)
COPY .air.toml .air.toml

# Command for development
CMD ["air"]

# -----------------------
# Stage 2: Build for production
# -----------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

# -----------------------
# Stage 3: Production Image
# -----------------------
FROM alpine:latest AS prod

WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080

CMD ["./app"]
