# Use the official Go image as a base
FROM golang:1.23.3

# Set the Current Working Directory inside the container
WORKDIR /app

# Install air for live reloading (optional for development)
RUN go install github.com/air-verse/air@latest

# Download dependencies
COPY go.mod ./

# Conditionally copy go.sum only if it exists
RUN if [ -f go.sum ]; then cp go.sum ./; fi

RUN go mod download

RUN go install github.com/go-delve/delve/cmd/dlv@latest

# ENTRYPOINT ["dlv", "debug", "--headless", "--listen=:40000", "--api-version=2", "--log", "--", "./main"]