# Use the official Golang image as the base image
FROM golang:1.17.8

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o /usr/local/bin/check ./main.go

# Expose the port that the application will run on
EXPOSE 3000

# Command to run the executable
CMD [ "/usr/local/bin/check" ]
