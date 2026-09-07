# Stage 1: build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o reminder-api ./cmd/api

# Stage 2: minimal runtime image
FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
WORKDIR /app
COPY --from=builder /app/reminder-api .
EXPOSE 8080
CMD ["./reminder-api"]