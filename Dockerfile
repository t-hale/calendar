FROM golang:1.24

WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary. Adjust the path "./cloud-run/create-calendar" if your main.go is elsewhere.
RUN go build -o main ./cloud-run

# Run the binary
CMD ["/app/main"]