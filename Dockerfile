# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Build the Go server
RUN go build -o server .

# Expose the port that the server will listen on
EXPOSE 8080

# Start the server when the container is run
CMD ["./main"]