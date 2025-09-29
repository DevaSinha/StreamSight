FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o worker .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/worker .

ENV API_ENDPOINT=http://api:8080

CMD ["./worker"]
