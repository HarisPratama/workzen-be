FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o main .

# ← install migrate binary di builder stage
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/database/migrations /app/database/migrations
# ← copy migrate binary dari builder
COPY --from=builder /go/bin/migrate /app/migrate

EXPOSE 8080
CMD ["/app/main"]