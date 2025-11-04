# syntax=docker/dockerfile:1
# Base builder
FROM golang:1.25.2-alpine AS build
ENV CGO_ENABLED=0
WORKDIR /app

# Install git for modules that use it and leverage Docker layer caching by
# downloading dependencies before copying the full source tree.
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Run tests and build the provider binary
RUN go test ./... && go build -o terraform-provider-validatefx .


# Final image
FROM hashicorp/terraform:1.9.8

# Copy provider binary for direct execution
COPY --from=build /app/terraform-provider-validatefx /usr/local/bin/terraform-provider-validatefx

# Install provider into Terraform override path
ENV TF_PLUGIN_DIR="/root/.terraform.d/plugins"
RUN mkdir -p ${TF_PLUGIN_DIR}/registry.terraform.io/the-devops-daily/validatefx/0.0.1/linux_amd64
COPY --from=build /app/terraform-provider-validatefx ${TF_PLUGIN_DIR}/registry.terraform.io/the-devops-daily/validatefx/0.0.1/linux_amd64/terraform-provider-validatefx_v0.0.1

# Provide workspace for Terraform configuration
WORKDIR /workspace

ENTRYPOINT []
CMD ["/bin/sh"]
