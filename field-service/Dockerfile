FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git openssh tzdata build-case

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

EXPOSE 8002

COPY --from=builder /app/field-service ./field-service

ENTRYPOINT["/app/field-service"]