# Start from a base image with Go installed
FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies 
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose port (if your app listens on a port)
EXPOSE 8088

# Command to run the executable
CMD ["./main"]
