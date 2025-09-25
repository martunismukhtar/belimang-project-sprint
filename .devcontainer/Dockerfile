FROM golang:1.25.0-alpine

ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Install dependencies (bash, etc.)
RUN apk add --no-cache bash git build-base curl

# Create group/user with default UID/GID (1000)
RUN addgroup -g 1000 megroup && \
    adduser -S -u 1000 -G megroup -s /bin/bash me

# Adjust UID/GID if build args differ from default
RUN if [ "$USER_UID" != "1000" ] || [ "$USER_GID" != "1000" ]; then \
      deluser me; \
      delgroup megroup; \
      addgroup -g $USER_GID megroup && \
      adduser -S -u $USER_UID -G megroup -s /bin/bash me; \
    fi

# Switch to non-root user
USER me

# Set GOPATH/bin to PATH so Air is accessible
ENV PATH=$PATH:/go/bin

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /workspace
