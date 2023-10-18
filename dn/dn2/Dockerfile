FROM golang:1.21 

# Set up a working directory
WORKDIR /app


# Copy the source code and generated Go files
COPY . .


# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify


# Compile the application
RUN go build -o main .

# Set the command to run the application
CMD ["./main"]
