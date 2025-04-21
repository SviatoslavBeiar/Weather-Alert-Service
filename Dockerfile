
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-alert-service ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/weather-alert-service .

COPY .env .env

EXPOSE 8080

CMD ["./weather-alert-service"]
