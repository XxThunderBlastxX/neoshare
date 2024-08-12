# Stage 1: Build the application
FROM golang:1.22.5-bullseye AS builder

# Install Bun
RUN apt install unzip
RUN curl -fsSL https://bun.sh/install | bash

# Install go-templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Install dependencies and build the application
RUN make build

# Stage 2: Create the runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/neoshare /app/neoshare

# Expose the port the app runs on
EXPOSE 8080

# Entrypoint command for the container
ENTRYPOINT ["/app/neoshare"]