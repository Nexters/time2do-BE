# Use the official Golang image as the base image
FROM golang:1.20

ENV SWAGGER_UI_CORS=true
ENV SWAGGER_UI_CORS_MAX_AGE=3600
ENV SWAGGER_UI_CORS_ALLOWED_ORIGINS=*
ENV SWAGGER_UI_CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
ENV SWAGGER_UI_CORS_ALLOWED_HEADERS=Content-Type,Authorization

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
