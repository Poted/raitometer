FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate .
COPY internal/migrations ./internal/migrations
COPY ./scripts/entrypoint.sh .

RUN chmod +x ./entrypoint.sh

EXPOSE 8080

CMD ["./entrypoint.sh"]