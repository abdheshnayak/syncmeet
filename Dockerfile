# Step 1: Build the Go program
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod tidy

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o myapp .  

# Step 2: Create the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port your app will run on (optional)
EXPOSE 8080

# Command to run the application
CMD ["/root/myapp"]
