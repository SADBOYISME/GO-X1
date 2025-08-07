# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# -o /app/main specifies the output file name
# CGO_ENABLED=0 disables Cgo to create a statically linked binary
# -ldflags="-s -w" strips debugging information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file (optional, but good for local docker-compose)
# You should ideally manage secrets more securely in production
COPY .env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
