FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./cmd/api

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/user-api /app/user-api

EXPOSE 8080
CMD ["./user-api"]

