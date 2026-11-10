# Stage 1: Build the Go application
FROM golang:1.27.4@sha256:08bcb99a38f85a265f9a93d1d3323e2440eec333380f3ea26ee0dec9264bb41c AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-backup-validator .

# Stage 2: Runtime image
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc

# renovate: datasource=github-tags depName=openbao/openbao
ARG VERSION_OPENBAO=2.6.2

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates jq wget procps && \
    rm -rf /var/lib/apt/lists/*

# Download and install Bao via the GlueOps Nexus raw proxy for GitHub (same one
# GlueKube uses). install_vault.sh uses the same proxy at runtime.
ADD https://repo.gpkg.io/repository/raw-github/openbao/openbao/releases/download/v${VERSION_OPENBAO}/openbao_${VERSION_OPENBAO}_linux_amd64.tar.gz /tmp/bao.tar.gz
RUN tar -xzf /tmp/bao.tar.gz bao && mv bao /usr/bin/bao && rm /tmp/bao.tar.gz

WORKDIR /app

# Copy binary and required files from builder
COPY --from=builder /app/vault-backup-validator .
COPY --from=builder /app/vault/configs ./vault/configs
COPY --from=builder /app/vault/scripts ./vault/scripts

EXPOSE 8080

CMD ["/app/vault-backup-validator"]
