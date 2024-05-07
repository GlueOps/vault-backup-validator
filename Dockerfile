# Stage 1: Build the Go application
FROM golang:1.21.10@sha256:392d2b634cba642c48e23b22949af823d42f4e722ca2d9f519133445e5a4cbba

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
RUN wget https://releases.hashicorp.com/vault/1.14.0/vault_1.14.0_linux_amd64.zip -O vault.zip && \
    unzip vault.zip -d /usr/local/bin && \
    rm vault.zip

EXPOSE 8080

# Start the application as root
CMD ["/app/vault-backup-validator"]
