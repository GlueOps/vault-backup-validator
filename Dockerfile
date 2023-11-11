# Stage 1: Build the Go application
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o vault-backup-validator .

RUN apt-get update
RUN apt-get install unzip -y
RUN apt-get install jq -y
#Download and install Vault
RUN wget https://releases.hashicorp.com/vault/1.15.0/vault_1.15.0_linux_amd64.zip -O vault.zip && \
    unzip vault.zip -d /usr/local/bin && \
    rm vault.zip

EXPOSE 8080

# Start the application as root
CMD ["/app/vault-backup-validator"]
