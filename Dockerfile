FROM golang:1.24 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main ./cloud-run/main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /app/main

# Make the binary executable
RUN chmod +x /app/main

# Service must listen to $PORT environment variable.
# This default value facilitates local development.
ENV PORT 8080

# Expose the port your application listens on (Cloud Run uses the PORT environment variable)
EXPOSE 0.0.0.0:$PORT:$PORT


# Run the web service on container startup.
CMD ["/app/main"]