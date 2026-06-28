# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy both modules (grpcsso depends on grpcsso-protos via replace directive)
COPY grpcsso-protos/ ./grpcsso-protos/
COPY grpcsso/ ./grpcsso/

WORKDIR /build/grpcsso

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/main.go

# Run stage
FROM alpine:3.22

RUN apk add --no-cache openssl

WORKDIR /app

COPY --from=builder /app .
COPY grpcsso/config.yaml .
COPY grpcsso/setup.sh .
COPY grpcsso/migrations/ ./migrations/

RUN chmod +x setup.sh && ./setup.sh

EXPOSE 3100

CMD ["./app"]
