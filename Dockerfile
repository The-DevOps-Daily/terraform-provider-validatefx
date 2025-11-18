# syntax=docker/dockerfile:1
# Multi-stage build for terraform-provider-validatefx

# Stage 1: Build
FROM golang:1.25.2-alpine AS build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Install dependencies (git for go modules)
RUN apk add --no-cache git ca-certificates

# Leverage Docker cache layers - copy go.mod/go.sum first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Run tests and build
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -v ./... && \
    go build -ldflags="-w -s" -trimpath -o terraform-provider-validatefx .

# Stage 2: Runtime
FROM hashicorp/terraform:1.9.8

# Copy CA certificates from build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy provider binary
COPY --from=build /build/terraform-provider-validatefx /usr/local/bin/terraform-provider-validatefx

# Install provider into Terraform plugin directory
ENV TF_PLUGIN_DIR="/root/.terraform.d/plugins" \
    PROVIDER_VERSION="0.0.1" \
    PROVIDER_NAMESPACE="registry.terraform.io/the-devops-daily/validatefx"

RUN mkdir -p ${TF_PLUGIN_DIR}/${PROVIDER_NAMESPACE}/${PROVIDER_VERSION}/linux_amd64 && \
    ln -s /usr/local/bin/terraform-provider-validatefx \
          ${TF_PLUGIN_DIR}/${PROVIDER_NAMESPACE}/${PROVIDER_VERSION}/linux_amd64/terraform-provider-validatefx_v${PROVIDER_VERSION}

# Set up workspace
WORKDIR /workspace

# Use exec form for proper signal handling
ENTRYPOINT []
CMD ["/bin/sh"]
