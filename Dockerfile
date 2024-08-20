# Stage 1: Build the Go application
FROM golang:1.22.6@sha256:367bb5295d3103981a86a572651d8297d6973f2ec8b62f716b007860e22cbc25

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o vault-backup-validator . && \
    apt-get update && \
    apt-get install -y --no-install-recommends unzip jq && \
    rm -rf /var/lib/apt/lists/*
    
#Download and install Vault
ADD https://releases.hashicorp.com/vault/1.14.0/vault_1.14.0_linux_amd64.zip /usr/local/bin/

EXPOSE 8080

# Start the application as root
CMD ["/app/vault-backup-validator"]