# Use the official Golang image as the base image
FROM golang:1.20

# Set the working directory
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Build the Go server
RUN go build -o server .

# Expose the port that the server will listen on
EXPOSE 8888

RUN chmod +x docker-entrypoint.sh wait-for-it.sh
ENTRYPOINT ./docker-entrypoint.sh
