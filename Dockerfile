# Stage 1: Build the Go application
FROM golang:1.25.3@sha256:8c945d3e25320e771326dafc6fb72ecae5f87b0f29328cbbd87c4dff506c9135

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
    

# renovate: datasource=github-tags depName=openbao/openbao
ARG VERSION_OPENBAO=2.4.3
ENV CACHED_OPENBAO_VERSION=${VERSION_OPENBAO}
  
#Download and install Bao
ADD https://github.com/openbao/openbao/releases/download/v${VERSION_OPENBAO}/bao_${VERSION_OPENBAO}_Linux_x86_64.tar.gz /tmp/bao_${VERSION_OPENBAO}_Linux_x86_64.tar.gz


# Unzip the Bao binary and clean up
RUN tar -xzvf /tmp/bao_${VERSION_OPENBAO}_Linux_x86_64.tar.gz bao && mv bao /usr/bin/bao && rm /tmp/bao_${VERSION_OPENBAO}_Linux_x86_64.tar.gz

EXPOSE 8080

# Start the application as root
CMD ["/app/vault-backup-validator"]
