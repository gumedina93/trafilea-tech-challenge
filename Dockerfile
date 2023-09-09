# Use the official Go image as the base image.
# Set the Go version and enable Go modules.
FROM golang:1.21 AS build

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files and download dependencies.
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application source code.
COPY . .

# Build the Go application with optimizations.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd

# Create a minimal runtime image to reduce size.
FROM scratch

# Copy only the built binary from the previous stage.
COPY --from=build /app/app /app

# Set the entry point for the container.
ENTRYPOINT ["/app"]