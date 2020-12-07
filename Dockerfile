# builder group for container image cloudrun-storage server binary
FROM golang:1.15-buster as builder
ARG VERS_TAG
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY ./cloudrun ./cloudrun/
RUN go build -ldflags="-X 'github.com/salzr/ww2020/cloudrun/pkg/version.Version=${VERS_TAG}'" -mod=readonly -v -o server ./cloudrun/cmd/cloudrun-storage/

# build light container and copy server binary from builder stage
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/server /app/server