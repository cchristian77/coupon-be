# Build Stage
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on

RUN apk add --no-cache git

ARG WORK_DIRECTORY="/usr/src/app"

WORKDIR ${WORK_DIRECTORY}

# Copy module files first (better cache)
COPY ./go.mod .
COPY ./go.sum .
COPY . .
RUN go mod download

# Copy source code
COPY . .

# Build binary from cmd/web
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/web ./cmd/web

# Runtime Stage
FROM alpine:latest

RUN adduser -D appuser
WORKDIR ${WORK_DIRECTORY}

# Copy binary
COPY --from=builder /bin/web .

# Copy env file (from project root)
COPY env.json .

USER appuser

# Run app
CMD ["./coupon_be"]