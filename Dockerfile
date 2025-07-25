# Dockerfile for Go Fiber + MongoDB backend (Railway deployment)
FROM golang:1.24.2-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY --from=builder /app/credentials.json* ./
COPY --from=builder /app/uploads ./uploads
EXPOSE 8080
CMD ["./app"]
