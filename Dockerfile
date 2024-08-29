# Use the official Golang image with Go 1.20
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application on port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
