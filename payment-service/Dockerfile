FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git openssh tzdata build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment-service .

FROM debian:bullseye-slim

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
    ca-certificates tzdata curl wkhtmltopdf fonts-dejavu fonts-droid-fallback fonts-freefont-ttf fonts-liberation && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    mkdir /app

WORKDIR /app

EXPOSE 8003

COPY --from=builder /app/payment-service ./payment-service
COPY --from=builder /app/template ./template

ENTRYPOINT ["/app/payment-service"]