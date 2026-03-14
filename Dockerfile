FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -ldflags="-s -w" -o main .

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/database/migrations /app/database/migrations

EXPOSE 8080

CMD ["/app/main"]
