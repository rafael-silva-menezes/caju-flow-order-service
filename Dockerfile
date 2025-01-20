# Use an official Golang image as a base
FROM golang:alpine

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o app .

# Command to run the application
CMD ["./app"]
