# Stage 1: Build
FROM golang:1.25.0-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o fiber-app ./src   # adjust path if needed

# Stage 2: Run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/fiber-app .
COPY --from=builder /app/docs ./docs
RUN apk add --no-cache ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

EXPOSE 7880
CMD ["./fiber-app"]
