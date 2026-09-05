# Stage 1: Build the high-performance Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency configuration files
COPY go.mod ./

# Copy raw source code scripts
COPY main.go ./

# Compile into a static, zero-dependency binary executable
RUN CGO_ENABLED=0 GOOS=linux go build -o go-etl .

# Stage 2: Create a secure, ultra-lightweight runtime container
FROM alpine:3.19

WORKDIR /root/

# Copy only the compiled binary from the builder environment
COPY --from=builder /app/go-etl .

# Expose the microservice web API port
EXPOSE 8080

# Launch the data engine microservice
CMD ["./go-etl"]
