FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git openssh tzdata build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --no-cache add tzdata curl && \
    mkdir /app

WORKDIR /app

EXPOSE 8001

COPY --from=builder /app/user-service ./user-service

ENTRYPOINT ["/app/user-service"]