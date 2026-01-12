# Build Stage
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on

RUN apk add --no-cache git

ARG WORK_DIRECTORY="/usr/src/app"
WORKDIR ${WORK_DIRECTORY}

# Copy module files first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o coupon_be ./cmd/web


# Runtime Stage
FROM alpine:latest

ARG WORK_DIRECTORY="/usr/src/app"
WORKDIR ${WORK_DIRECTORY}

RUN adduser -D appuser

# Copy binary
COPY --from=builder /usr/src/app/coupon_be .

# Copy env file
COPY env.json .

USER appuser

CMD ["./coupon_be"]