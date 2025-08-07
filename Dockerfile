# Use official Golang image as base
FROM golang:1.24-alpine

# Working Directory
WORKDIR /PASSMAN

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port your app uses
EXPOSE 1019

# Command to run the app
CMD ["./main"]