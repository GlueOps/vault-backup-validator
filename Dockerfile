# Stage 1: Build the Go application
FROM golang:1.24.2@sha256:d9db32125db0c3a680cfb7a1afcaefb89c898a075ec148fdc2f0f646cc2ed509

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
