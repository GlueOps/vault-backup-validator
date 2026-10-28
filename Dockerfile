# Stage 1: Build the Go application
FROM golang:1.27.3@sha256:b76020929b709ddb6202afc05410a0faefae20ae05e27ad548b3870092688a69 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-backup-validator .

# Stage 2: Runtime image
FROM debian:bookworm-slim@sha256:99c0b25c46d0214d5a2091092e7908907dcaecd6a962cde0ac700acee222b6f3

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
