# Stage 1: Build the Go application
FROM golang:1.25.4@sha256:698183780de28062f4ef46f82a79ec0ae69d2d22f7b160cf69f71ea8d98bf25d AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-backup-validator .

# Stage 2: Runtime image
FROM debian:bookworm-slim

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
