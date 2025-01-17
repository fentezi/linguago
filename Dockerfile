FROM golang:1.23.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main ./cmd/translator/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main .
COPY --from=builder /build/config.yml config.yml
COPY --from=builder /build/audio /app/audio
COPY .env .env

EXPOSE 8080

CMD ["./main", "--config=./config.yml"]
