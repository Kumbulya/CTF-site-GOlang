FROM golang:1.21.4-alpine
# Set the working directory in the container
WORKDIR /app
COPY . .
# Build your Go application
RUN go mod download
RUN go build -o application ./web 
# Expose port 1337
EXPOSE 1337
# Run the command to start your application
CMD ["./application"]