# Stage 1: Build the Go application
FROM golang:1.23.4 as builder
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod ./
RUN go mod tidy

# Copy the main.go file and other application files
COPY . ./

# Build the main.go file into a binary named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Stage 2: Create a minimal image using distroless
FROM gcr.io/distroless/static
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Command to run the application
CMD ["/app/main"]