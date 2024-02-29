# Use the official Golang image as the base image
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /check_point

# Copy the local package files to the container's workspace
COPY . /check_point

# Build the Go application
# RUN go mod download
RUN go build -o check-point ./main.go

EXPOSE 8080:8080

CMD [ "./check-point" ]
