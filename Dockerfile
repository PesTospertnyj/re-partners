FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/pack-api ./cmd/api

FROM alpine:3.20

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/pack-api /app/pack-api
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/pack-api"]
