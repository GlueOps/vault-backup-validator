# Stage 1: Build the Go application
FROM golang:1.25.1@sha256:8305f5fa8ea63c7b5bc85bd223ccc62941f852318ebfbd22f53bbd0b358c07e1

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
ADD https://releases.hashicorp.com/vault/1.14.0/vault_1.14.0_linux_amd64.zip /tmp/vault.zip

# Unzip the Vault binary and clean up
RUN unzip /tmp/vault.zip -d /usr/local/bin/ && \
    rm /tmp/vault.zip

EXPOSE 8080

# Start the application as root
CMD ["/app/vault-backup-validator"]
