# Stage 1: Build the Go application
FROM golang:1.26.0@sha256:c83e68f3ebb6943a2904fa66348867d108119890a2c6a2e6f07b38d0eb6c25c5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-backup-validator .

# Stage 2: Runtime image
FROM debian:bookworm-slim@sha256:98f4b71de414932439ac6ac690d7060df1f27161073c5036a7553723881bffbe

# renovate: datasource=github-tags depName=openbao/openbao
ARG VERSION_OPENBAO=2.4.4
ENV CACHED_OPENBAO_VERSION=${VERSION_OPENBAO}

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates jq wget procps && \
    rm -rf /var/lib/apt/lists/*

# Download and install Bao
ADD https://github.com/openbao/openbao/releases/download/v${VERSION_OPENBAO}/bao_${VERSION_OPENBAO}_Linux_x86_64.tar.gz /tmp/bao.tar.gz
RUN tar -xzf /tmp/bao.tar.gz bao && mv bao /usr/bin/bao && rm /tmp/bao.tar.gz

WORKDIR /app

# Copy binary and required files from builder
COPY --from=builder /app/vault-backup-validator .
COPY --from=builder /app/vault/configs ./vault/configs
COPY --from=builder /app/vault/scripts ./vault/scripts

EXPOSE 8080

CMD ["/app/vault-backup-validator"]
