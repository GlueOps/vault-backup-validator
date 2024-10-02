# Stage 1: Build the Go application
FROM golang:1.22.8@sha256:628529a29f130a8ab336b994be99d134ce98cd23b8f2052d8995678681e97ca2

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
