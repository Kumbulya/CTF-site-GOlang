FROM golang:1.16-alpine

# Set the working directory in the container
WORKDIR /app

COPY . .
# Build your Go application
RUN go build main.go

# Expose port 1337
EXPOSE 1337

# Run the command to start your application
CMD ["./main"]